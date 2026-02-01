package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"spotiflac/backend"

	"github.com/spf13/cobra"
)

var (
	outputDir            string
	service              string
	quality              string
	audioFormat          string
	filenameFormat       string
	includeTrackNumber   bool
	embedLyrics          bool
	embedMaxQualityCover bool
	allowFallback        bool
	autoDownload         bool
	autoOrder            string
	autoQuality          string
	batchDelay           float64
	region               string
	quiet                bool
	verbose              bool
	dumpJSON             bool
	writeInfoJSON        bool
	version              = "1.0.0"
)

var rootCmd = &cobra.Command{
	Use:   "spotiflac [URL]",
	Short: "SpotiFLAC - Download Spotify tracks in lossless FLAC format",
	Long: `SpotiFLAC - Get Spotify tracks in true FLAC from Tidal, Qobuz & Amazon Music

Examples:
  spotiflac https://open.spotify.com/track/...
  spotiflac https://open.spotify.com/album/... --service tidal
  spotiflac https://open.spotify.com/playlist/... --auto --output ./Music
  spotiflac https://open.spotify.com/track/... --quality HI_RES --embed-lyrics
  
Services: tidal, qobuz, amazon, auto (tries multiple services)
`,
	Version: version,
	Args:    cobra.MinimumNArgs(1),
	Run:     runDownload,
}

func init() {
	// Output options
	rootCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory")
	rootCmd.Flags().StringVarP(&filenameFormat, "filename-format", "f", "title-artist", "Filename format: title-artist, artist-title, title, or custom template")
	rootCmd.Flags().BoolVarP(&includeTrackNumber, "track-number", "n", false, "Include track number in filename")

	// Service options
	rootCmd.Flags().StringVarP(&service, "service", "s", "auto", "Download service: tidal, qobuz, amazon, auto")
	rootCmd.Flags().BoolVarP(&autoDownload, "auto", "a", false, "Auto mode: try multiple services (same as --service auto)")
	rootCmd.Flags().StringVar(&autoOrder, "auto-order", "tidal-amazon-qobuz", "Service order for auto mode")
	rootCmd.Flags().StringVar(&autoQuality, "auto-quality", "24", "Auto quality: 16 or 24 bit")

	// Quality options
	rootCmd.Flags().StringVarP(&quality, "quality", "q", "", "Quality: LOSSLESS, HI_RES (Tidal) or 6, 7, 27 (Qobuz)")
	rootCmd.Flags().StringVar(&audioFormat, "format", "", "Audio format (alias for quality)")
	rootCmd.Flags().BoolVar(&allowFallback, "fallback", true, "Allow fallback to lower quality")

	// Metadata options
	rootCmd.Flags().BoolVarP(&embedLyrics, "embed-lyrics", "l", false, "Embed lyrics in file")
	rootCmd.Flags().BoolVarP(&embedMaxQualityCover, "max-quality-cover", "c", false, "Embed maximum quality album cover")

	// Batch options
	rootCmd.Flags().Float64VarP(&batchDelay, "delay", "d", 1.0, "Delay between downloads (seconds)")
	rootCmd.Flags().StringVarP(&region, "region", "r", "US", "Region code for matching")

	// Output options
	rootCmd.Flags().BoolVarP(&quiet, "quiet", "", false, "Quiet mode (minimal output)")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	// Metadata options
	rootCmd.Flags().BoolVarP(&dumpJSON, "dump-json", "j", false, "Print metadata as JSON and exit (no download)")
	rootCmd.Flags().BoolVar(&writeInfoJSON, "write-info-json", false, "Write metadata to .info.json file")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runDownload(cmd *cobra.Command, args []string) {
	spotifyURL := args[0]

	if autoDownload {
		service = "auto"
	}

	// Set quality based on service
	if audioFormat != "" {
		quality = audioFormat
	}

	if quality == "" {
		switch service {
		case "tidal":
			quality = "LOSSLESS"
		case "qobuz":
			if autoQuality == "24" {
				quality = "7"
			} else {
				quality = "6"
			}
		default:
			quality = "LOSSLESS"
		}
	}

	if !dumpJSON {
		logInfo("SpotiFLAC v%s - Starting download...", version)
		logInfo("URL: %s", spotifyURL)
	}

	// Initialize history DB (optional for CLI)
	if err := backend.InitHistoryDB("SpotiFLAC-CLI"); err != nil {
		logDebug("Warning: Failed to init history DB: %v", err)
	}
	defer backend.CloseHistoryDB()

	// Fetch Spotify metadata
	if !dumpJSON {
		logInfo("Fetching metadata from Spotify...")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	metadata, err := backend.GetFilteredSpotifyData(ctx, spotifyURL, true, time.Duration(batchDelay)*time.Second)
	if err != nil {
		logError("Failed to fetch metadata: %v", err)
		os.Exit(1)
	}

	// Dump JSON mode - print metadata and exit
	if dumpJSON {
		printMetadataJSONWithLinks(metadata, spotifyURL)
		return
	}

	// Write info JSON if requested
	if writeInfoJSON {
		writeMetadataJSON(metadata, spotifyURL)
	}

	// Determine content type and download based on actual return types
	switch data := metadata.(type) {
	case backend.TrackResponse:
		// Single track - convert to map format for compatibility
		trackMap := map[string]interface{}{
			"spotify_id":   data.Track.SpotifyID,
			"name":         data.Track.Name,
			"artists":      data.Track.Artists,
			"album_name":   data.Track.AlbumName,
			"album_artist": data.Track.AlbumArtist,
			"release_date": data.Track.ReleaseDate,
			"images":       data.Track.Images,
			"isrc":         data.Track.ISRC,
			"duration_ms":  float64(data.Track.DurationMS),
			"track_number": float64(data.Track.TrackNumber),
			"disc_number":  float64(data.Track.DiscNumber),
			"total_tracks": float64(data.Track.TotalTracks),
			"total_discs":  float64(data.Track.TotalDiscs),
			"copyright":    data.Track.Copyright,
			"publisher":    data.Track.Publisher,
		}
		downloadTrack(trackMap)

	case *backend.AlbumResponsePayload:
		// Album
		logInfo("Album: %s (%d tracks)", data.AlbumInfo.Name, len(data.TrackList))
		trackList := make([]interface{}, len(data.TrackList))
		for i, t := range data.TrackList {
			trackList[i] = map[string]interface{}{
				"spotify_id":   t.SpotifyID,
				"name":         t.Name,
				"artists":      t.Artists,
				"album_name":   t.AlbumName,
				"album_artist": t.AlbumArtist,
				"release_date": t.ReleaseDate,
				"images":       t.Images,
				"isrc":         t.ISRC,
				"track_number": float64(t.TrackNumber),
				"disc_number":  float64(t.DiscNumber),
				"total_tracks": float64(t.TotalTracks),
				"total_discs":  float64(t.TotalDiscs),
			}
		}
		downloadBatch(trackList, "album")

	case *backend.PlaylistResponsePayload:
		// Playlist
		logInfo("Playlist: %d tracks", len(data.TrackList))
		trackList := make([]interface{}, len(data.TrackList))
		for i, t := range data.TrackList {
			trackList[i] = map[string]interface{}{
				"spotify_id":   t.SpotifyID,
				"name":         t.Name,
				"artists":      t.Artists,
				"album_name":   t.AlbumName,
				"album_artist": t.AlbumArtist,
				"release_date": t.ReleaseDate,
				"images":       t.Images,
				"isrc":         t.ISRC,
				"track_number": float64(t.TrackNumber),
				"disc_number":  float64(t.DiscNumber),
				"total_tracks": float64(t.TotalTracks),
				"total_discs":  float64(t.TotalDiscs),
			}
		}
		downloadBatch(trackList, "playlist")

	case *backend.ArtistDiscographyPayload:
		// Artist discography
		logInfo("Artist: %s (%d tracks)", data.ArtistInfo.Name, len(data.TrackList))
		trackList := make([]interface{}, len(data.TrackList))
		for i, t := range data.TrackList {
			trackList[i] = map[string]interface{}{
				"spotify_id":   t.SpotifyID,
				"name":         t.Name,
				"artists":      t.Artists,
				"album_name":   t.AlbumName,
				"album_artist": t.AlbumArtist,
				"release_date": t.ReleaseDate,
				"images":       t.Images,
				"isrc":         t.ISRC,
				"track_number": float64(t.TrackNumber),
				"disc_number":  float64(t.DiscNumber),
				"total_tracks": float64(t.TotalTracks),
				"total_discs":  float64(t.TotalDiscs),
			}
		}
		downloadBatch(trackList, "artist")

	default:
		logError("Unknown content type: %T", metadata)
		os.Exit(1)
	}

	logSuccess("All downloads complete!")
}

func printMetadataJSON(metadata interface{}) {
	jsonData, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		logError("Failed to marshal metadata to JSON: %v", err)
		os.Exit(1)
	}
	fmt.Println(string(jsonData))
}

func printMetadataJSONWithLinks(metadata interface{}, spotifyURL string) {
	// Extract Spotify ID
	var spotifyID string
	switch data := metadata.(type) {
	case backend.TrackResponse:
		spotifyID = data.Track.SpotifyID
	case *backend.AlbumResponsePayload:
		if len(data.TrackList) > 0 {
			spotifyID = data.TrackList[0].SpotifyID
		}
	case *backend.PlaylistResponsePayload:
		if len(data.TrackList) > 0 {
			spotifyID = data.TrackList[0].SpotifyID
		}
	case *backend.ArtistDiscographyPayload:
		if len(data.TrackList) > 0 {
			spotifyID = data.TrackList[0].SpotifyID
		}
	}

	// Create output structure
	output := make(map[string]interface{})
	output["metadata"] = metadata

	// Get streaming URLs if we have a Spotify ID
	if spotifyID != "" {
		client := backend.NewSongLinkClient()
		urls, err := client.GetAllURLsFromSpotify(spotifyID, region)
		if err == nil {
			downloadLinks := make(map[string]interface{})

			if urls.TidalURL != "" {
				downloadLinks["tidal"] = map[string]string{
					"url":     urls.TidalURL,
					"quality": "LOSSLESS (16-bit) / HI_RES (24-bit)",
				}
			}
			if urls.AmazonURL != "" {
				downloadLinks["amazon"] = map[string]string{
					"url":     urls.AmazonURL,
					"quality": "Variable (auto-select)",
				}
			}

			// Add Qobuz note (requires ISRC lookup)
			downloadLinks["qobuz"] = map[string]string{
				"note":    "Available via ISRC lookup (requires Deezer)",
				"quality": "6 (16-bit) / 7 (24-bit) / 27 (Hi-Res)",
			}

			output["download_links"] = downloadLinks
			output["available_services"] = []string{}
			if urls.TidalURL != "" {
				output["available_services"] = append(output["available_services"].([]string), "tidal")
			}
			if urls.AmazonURL != "" {
				output["available_services"] = append(output["available_services"].([]string), "amazon")
			}
			output["available_services"] = append(output["available_services"].([]string), "qobuz")
		}
	}

	output["spotiflac_cli"] = map[string]interface{}{
		"version": version,
		"usage": map[string]string{
			"tidal":  fmt.Sprintf("spotiflac %s --service tidal --quality LOSSLESS", spotifyURL),
			"amazon": fmt.Sprintf("spotiflac %s --service amazon", spotifyURL),
			"qobuz":  fmt.Sprintf("spotiflac %s --service qobuz --quality 7", spotifyURL),
			"auto":   fmt.Sprintf("spotiflac %s --auto --auto-quality 24", spotifyURL),
		},
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		logError("Failed to marshal metadata to JSON: %v", err)
		os.Exit(1)
	}
	fmt.Println(string(jsonData))
}

func writeMetadataJSON(metadata interface{}, spotifyURL string) {
	// Extract filename from metadata
	var filename string
	switch data := metadata.(type) {
	case backend.TrackResponse:
		filename = sanitizeFilename(data.Track.Name + " - " + data.Track.Artists)
	case *backend.AlbumResponsePayload:
		filename = sanitizeFilename(data.AlbumInfo.Name)
	case *backend.PlaylistResponsePayload:
		filename = "playlist"
	case *backend.ArtistDiscographyPayload:
		filename = sanitizeFilename(data.ArtistInfo.Name)
	default:
		filename = "metadata"
	}

	jsonFilePath := filepath.Join(outputDir, filename+".info.json")
	jsonData, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		logError("Failed to marshal metadata to JSON: %v", err)
		return
	}

	if err := os.WriteFile(jsonFilePath, jsonData, 0644); err != nil {
		logError("Failed to write info JSON: %v", err)
		return
	}

	logInfo("Metadata saved to: %s", jsonFilePath)
}

func sanitizeFilename(name string) string {
	// Remove invalid filename characters
	invalid := []string{"<", ">", ":", "\"", "/", "\\", "|", "?", "*"}
	result := name
	for _, char := range invalid {
		result = strings.ReplaceAll(result, char, "_")
	}
	return result
}

func downloadTrack(trackData map[string]interface{}) {
	spotifyID, _ := trackData["spotify_id"].(string)
	trackName, _ := trackData["name"].(string)
	artistName, _ := trackData["artists"].(string)
	albumName, _ := trackData["album_name"].(string)
	albumArtist, _ := trackData["album_artist"].(string)
	releaseDate, _ := trackData["release_date"].(string)
	coverURL, _ := trackData["images"].(string)
	isrc, _ := trackData["isrc"].(string)
	durationMS, _ := trackData["duration_ms"].(float64)
	trackNumber, _ := trackData["track_number"].(float64)
	discNumber, _ := trackData["disc_number"].(float64)
	totalTracks, _ := trackData["total_tracks"].(float64)
	totalDiscs, _ := trackData["total_discs"].(float64)
	copyright, _ := trackData["copyright"].(string)
	publisher, _ := trackData["publisher"].(string)

	logInfo("\nDownloading: %s - %s", artistName, trackName)
	logDebug("ISRC: %s | Album: %s", isrc, albumName)

	var filePath string
	var downloadErr error

	if service == "auto" {
		filePath, downloadErr = downloadWithAutoFallback(
			spotifyID, isrc, trackName, artistName, albumName, albumArtist,
			releaseDate, coverURL, int(durationMS), int(trackNumber), int(discNumber),
			int(totalTracks), int(totalDiscs), copyright, publisher,
		)
	} else {
		filePath, downloadErr = downloadWithService(
			service, spotifyID, isrc, trackName, artistName, albumName, albumArtist,
			releaseDate, coverURL, int(trackNumber), int(discNumber),
			int(totalTracks), int(totalDiscs), copyright, publisher,
		)
	}

	if downloadErr != nil {
		logError("Download failed: %v", downloadErr)
		os.Exit(1)
	}

	if strings.HasPrefix(filePath, "EXISTS:") {
		logWarning("File already exists: %s", strings.TrimPrefix(filePath, "EXISTS:"))
	} else {
		logSuccess("Downloaded: %s", filePath)
	}
}

func downloadBatch(trackList []interface{}, contentType string) {
	total := len(trackList)
	succeeded := 0
	failed := 0
	skipped := 0

	for i, trackItem := range trackList {
		trackData := trackItem.(map[string]interface{})

		logInfo("\n[%d/%d] Processing track...", i+1, total)

		spotifyID, _ := trackData["spotify_id"].(string)
		trackName, _ := trackData["name"].(string)
		artistName, _ := trackData["artists"].(string)
		albumName, _ := trackData["album_name"].(string)
		albumArtist, _ := trackData["album_artist"].(string)
		releaseDate, _ := trackData["release_date"].(string)
		coverURL, _ := trackData["images"].(string)
		isrc, _ := trackData["isrc"].(string)
		trackNumber, _ := trackData["track_number"].(float64)
		discNumber, _ := trackData["disc_number"].(float64)
		totalTracks, _ := trackData["total_tracks"].(float64)
		totalDiscs, _ := trackData["total_discs"].(float64)
		copyright, _ := trackData["copyright"].(string)
		publisher, _ := trackData["publisher"].(string)

		var filePath string
		var downloadErr error

		if service == "auto" {
			filePath, downloadErr = downloadWithAutoFallback(
				spotifyID, isrc, trackName, artistName, albumName, albumArtist,
				releaseDate, coverURL, 0, int(trackNumber), int(discNumber),
				int(totalTracks), int(totalDiscs), copyright, publisher,
			)
		} else {
			filePath, downloadErr = downloadWithService(
				service, spotifyID, isrc, trackName, artistName, albumName, albumArtist,
				releaseDate, coverURL, int(trackNumber), int(discNumber),
				int(totalTracks), int(totalDiscs), copyright, publisher,
			)
		}

		if downloadErr != nil {
			logError("Failed: %s - %s (%v)", artistName, trackName, downloadErr)
			failed++
		} else if strings.HasPrefix(filePath, "EXISTS:") {
			logWarning("Skipped (exists): %s - %s", artistName, trackName)
			skipped++
		} else {
			logSuccess("Downloaded: %s - %s", artistName, trackName)
			succeeded++
		}

		// Delay between downloads
		if i < total-1 && batchDelay > 0 {
			time.Sleep(time.Duration(batchDelay * float64(time.Second)))
		}
	}

	logInfo("\n=== Summary ===")
	logInfo("Total: %d | Success: %d | Failed: %d | Skipped: %d", total, succeeded, failed, skipped)
}

func downloadWithAutoFallback(spotifyID, isrc, trackName, artistName, albumName, albumArtist, releaseDate, coverURL string, durationMS, trackNumber, discNumber, totalTracks, totalDiscs int, copyright, publisher string) (string, error) {
	// Get streaming URLs
	client := backend.NewSongLinkClient()
	urls, err := client.GetAllURLsFromSpotify(spotifyID, region)
	if err != nil {
		return "", fmt.Errorf("failed to get streaming URLs: %w", err)
	}

	// Parse auto order
	orderServices := strings.Split(autoOrder, "-")

	// Determine qualities
	is24Bit := autoQuality == "24"
	tidalQuality := "LOSSLESS"
	qobuzQuality := "6"
	if is24Bit {
		tidalQuality = "HI_RES_LOSSLESS"
		qobuzQuality = "7"
	}

	var lastErr error

	for _, svc := range orderServices {
		svc = strings.TrimSpace(svc)

		switch svc {
		case "tidal":
			if urls.TidalURL != "" {
				logDebug("Trying Tidal...")
				downloader := backend.NewTidalDownloader("")
				filePath, err := downloader.DownloadByURLWithFallback(
					urls.TidalURL, outputDir, tidalQuality, filenameFormat,
					includeTrackNumber, trackNumber, trackName, artistName,
					albumName, albumArtist, releaseDate, true, coverURL,
					embedMaxQualityCover, trackNumber, discNumber, totalTracks,
					totalDiscs, copyright, publisher,
					fmt.Sprintf("https://open.spotify.com/track/%s", spotifyID),
					allowFallback,
				)
				if err == nil {
					return filePath, nil
				}
				lastErr = err
				logDebug("Tidal failed, trying next service...")
			}

		case "amazon":
			if urls.AmazonURL != "" {
				logDebug("Trying Amazon Music...")
				downloader := backend.NewAmazonDownloader()
				filePath, err := downloader.DownloadByURL(
					urls.AmazonURL, outputDir, quality, filenameFormat,
					"", "", includeTrackNumber, trackNumber, trackName,
					artistName, albumName, albumArtist, releaseDate, coverURL,
					trackNumber, discNumber, totalTracks, embedMaxQualityCover,
					totalDiscs, copyright, publisher,
					fmt.Sprintf("https://open.spotify.com/track/%s", spotifyID),
				)
				if err == nil {
					return filePath, nil
				}
				lastErr = err
				logDebug("Amazon failed, trying next service...")
			}

		case "qobuz":
			logDebug("Trying Qobuz...")
			// Get Deezer ISRC for Qobuz
			songlinkClient := backend.NewSongLinkClient()
			deezerURL, err := songlinkClient.GetDeezerURLFromSpotify(spotifyID)
			if err == nil {
				deezerISRC, err := backend.GetDeezerISRC(deezerURL)
				if err == nil && deezerISRC != "" {
					downloader := backend.NewQobuzDownloader()
					filePath, err := downloader.DownloadByISRC(
						deezerISRC, outputDir, qobuzQuality, filenameFormat,
						includeTrackNumber, trackNumber, trackName, artistName,
						albumName, albumArtist, releaseDate, true, coverURL,
						embedMaxQualityCover, trackNumber, discNumber, totalTracks,
						totalDiscs, copyright, publisher,
						fmt.Sprintf("https://open.spotify.com/track/%s", spotifyID),
						allowFallback,
					)
					if err == nil {
						return filePath, nil
					}
					lastErr = err
				}
			}
			logDebug("Qobuz failed, trying next service...")
		}
	}

	if lastErr != nil {
		return "", fmt.Errorf("all services failed: %w", lastErr)
	}
	return "", fmt.Errorf("no services available")
}

func downloadWithService(svc, spotifyID, isrc, trackName, artistName, albumName, albumArtist, releaseDate, coverURL string, trackNumber, discNumber, totalTracks, totalDiscs int, copyright, publisher string) (string, error) {
	spotifyURL := fmt.Sprintf("https://open.spotify.com/track/%s", spotifyID)

	switch svc {
	case "tidal":
		downloader := backend.NewTidalDownloader("")
		return downloader.Download(
			spotifyID, outputDir, quality, filenameFormat,
			includeTrackNumber, trackNumber, trackName, artistName,
			albumName, albumArtist, releaseDate, true, coverURL,
			embedMaxQualityCover, trackNumber, discNumber, totalTracks,
			totalDiscs, copyright, publisher, spotifyURL, allowFallback,
		)

	case "qobuz":
		// Get Deezer ISRC
		client := backend.NewSongLinkClient()
		deezerURL, err := client.GetDeezerURLFromSpotify(spotifyID)
		if err != nil {
			return "", fmt.Errorf("failed to get Deezer URL: %w", err)
		}

		deezerISRC, err := backend.GetDeezerISRC(deezerURL)
		if err != nil {
			return "", fmt.Errorf("failed to get ISRC: %w", err)
		}

		downloader := backend.NewQobuzDownloader()
		return downloader.DownloadByISRC(
			deezerISRC, outputDir, quality, filenameFormat,
			includeTrackNumber, trackNumber, trackName, artistName,
			albumName, albumArtist, releaseDate, true, coverURL,
			embedMaxQualityCover, trackNumber, discNumber, totalTracks,
			totalDiscs, copyright, publisher, spotifyURL, allowFallback,
		)

	case "amazon":
		downloader := backend.NewAmazonDownloader()
		return downloader.DownloadBySpotifyID(
			spotifyID, outputDir, quality, filenameFormat,
			"", "", includeTrackNumber, trackNumber, trackName,
			artistName, albumName, albumArtist, releaseDate, coverURL,
			trackNumber, discNumber, totalTracks, embedMaxQualityCover,
			totalDiscs, copyright, publisher, spotifyURL,
		)

	default:
		return "", fmt.Errorf("unknown service: %s", svc)
	}
}

// Logging functions
func logInfo(format string, args ...interface{}) {
	if !quiet {
		fmt.Printf("[INFO] "+format+"\n", args...)
	}
}

func logSuccess(format string, args ...interface{}) {
	if !quiet {
		fmt.Printf("✓ [SUCCESS] "+format+"\n", args...)
	}
}

func logError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "✗ [ERROR] "+format+"\n", args...)
}

func logWarning(format string, args ...interface{}) {
	if !quiet {
		fmt.Printf("⚠ [WARNING] "+format+"\n", args...)
	}
}

func logDebug(format string, args ...interface{}) {
	if verbose {
		fmt.Printf("[DEBUG] "+format+"\n", args...)
	}
}
