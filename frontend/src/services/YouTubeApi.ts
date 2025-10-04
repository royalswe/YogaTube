// Extend the Window interface to include YT
declare global {
    interface Window {
        YT: any;
        onYouTubeIframeAPIReady?: () => void;
    }
}

export class YouTubeApi {
    private player: any = null;
    private ready: boolean = false;

    constructor(playerId: string, videoId: string, onReady?: () => void) {
        // Load YouTube IFrame API if not already loaded
        if (!window.YT) {
            const tag = document.createElement('script');
            tag.src = 'https://www.youtube.com/iframe_api';
            document.body.appendChild(tag);
        }
        window.onYouTubeIframeAPIReady = () => {
            this.player = new window.YT.Player(playerId, {
                videoId,
                events: {
                    onReady: () => {
                        this.ready = true;
                        if (onReady) onReady();
                    }
                }
            });
        };
    }

    loadVideoById(videoId: string) {
        if (this.ready && this.player) {
            this.player.cueVideoById(videoId);
        }
    }

    togglePlayPause() {
        if (this.ready) {
            if (this.player?.getPlayerState() === 1) { // Not playing
                this.player?.pauseVideo();
            } else {
                this.player?.playVideo();
            }
        }
    }

    toggleFullscreen() {
        if (this.ready) {
            if (document.fullscreenElement) {
                document.exitFullscreen();
            } else {
                this.player?.getIframe().requestFullscreen();
            }
        }
    }
}
