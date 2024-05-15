export class VideoPlayer {
    volumeIncrement = 0.05;     // 5%
    videoScrubIncrement = 5;    // seconds

    player;

    buttons;


    constructor() {
        let player = document.getElementsByTagName('video');
        if (player.length === 0) {
            return;
        }
        this.player = player[0];

        this.#LoadButtons();
        this.#AddListeners();
        this.#UpdateUI();

    }

    #LoadButtons() {
        let play = document.getElementById('play');
        let rewind = document.getElementById('rewind');
        let forward = document.getElementById('forward');
        let mute = document.getElementById('mute');
        let fullScreen = document.getElementById('full-screen');

        this.buttons = {
            play: play,
            rewind: rewind,
            forward: forward,
            mute: mute,
            fullScreen: fullScreen,
        }

    }

    #AddListeners() {
        // buton listeners
        this.player.addEventListener('click', this.Play.bind(this));
        this.buttons.play.addEventListener('click', this.Play.bind(this));
        this.buttons.rewind.addEventListener('click', this.Rewind.bind(this));
        this.buttons.forward.addEventListener('click', this.Forward.bind(this));
        this.buttons.mute.addEventListener('click', this.Mute.bind(this));
        this.buttons.fullScreen.addEventListener('click', this.ToggleFullScreen.bind(this));

        // keyboard listeners
        document.addEventListener('keydown', this.#KeyboardListener.bind(this));

        // video listeners
        //   - play
        this.player.addEventListener('play', function () {
            this.UpdatePlayButton("pause")
            document.getElementById('progress-control').max = this.player.duration;
        }.bind(this));

        this.player.addEventListener('pause', function () {
            this.UpdatePlayButton("play_arrow")
        }.bind(this));

        this.player.addEventListener('ended', function () {
            this.UpdatePlayButton("replay")
        }.bind(this));
        //   - seeked
        this.player.addEventListener('seeked', this.UpdateEndsAt.bind(this));
        //   - volume
        let volumeControl = document.getElementById('volume-control');
        volumeControl.addEventListener('input', function (e) {
            this.player.volume = e.target.value;
            this.UpdateVolumeControl();
        }.bind(this));

        //   - progress
        this.player.addEventListener('timeupdate', this.UpdateCurrentProgress.bind(this));
        let progressControl = document.getElementById('progress-control');
        progressControl.max = this.player.duration;
        progressControl.addEventListener('input', function (e) {
            this.player.currentTime = e.target.value;
        }.bind(this));
    }

    #KeyboardListener(e) {
        switch (e.code) {
            case 'Space':
                this.Play();
                break;
            case 'ArrowLeft':
                this.Rewind();
                break;
            case 'ArrowRight':
                this.Forward();
                break;
            case 'Home':
                this.player.currentTime = 0;
                break;
            case 'End':
                this.player.currentTime = this.player.duration;
                break;
            case 'ArrowUp':
                this.VolumeUp();
                break;
            case 'ArrowDown':
                this.VolumeDown();
                break;
            case 'Digit1':
                this.player.currentTime = this.player.duration * 0.1;
                break;
            case 'Digit2':
                this.player.currentTime = this.player.duration * 0.2;
                break;
            case 'Digit3':
                this.player.currentTime = this.player.duration * 0.3;
                break;
            case 'Digit4':
                this.player.currentTime = this.player.duration * 0.4;
                break;
            case 'Digit5':
                this.player.currentTime = this.player.duration * 0.5;
                break;
            case 'Digit6':
                this.player.currentTime = this.player.duration * 0.6;
                break;
            case 'Digit7':
                this.player.currentTime = this.player.duration * 0.7;
                break;
            case 'Digit8':
                this.player.currentTime = this.player.duration * 0.8;
                break;
            case 'Digit9':
                this.player.currentTime = this.player.duration * 0.9;
                break;
            case 'KeyF':
                this.ToggleFullScreen();
                break;
            default:
                break;
        }
    }

    #UpdateUI() {
        this.UpdatePlayButton();
        this.UpdateEndsAt();
        this.UpdateCurrentProgress();
        this.UpdateVolumeControl();
    }

    Play() {
        this.player.paused ? this.player.play() : this.player.pause();
        this.UpdateEndsAt();
    }

    UpdatePlayButton(newButtonText) {
        let buttontext = this.buttons.play.children;
        if (!newButtonText) newButtonText = this.player.paused ? "play_arrow" : "pause";
        if (buttontext.length > 0) {
            buttontext[0].innerHTML = newButtonText;
        } else {
            this.buttons.play.innerHTML = newButtonText;
        }
    }

    Rewind() {
        this.player.currentTime -= this.videoScrubIncrement;
        this.UpdateEndsAt();
    }

    Forward() {
        this.player.currentTime += this.videoScrubIncrement;
        this.UpdateEndsAt();
    }


    UpdateEndsAt() {
        let endsAt = document.getElementById('ends-at');
        let currentTime = Date.now();
        currentTime += (this.player.duration - this.player.currentTime) * 1000;

        endsAt.innerHTML = new Date(currentTime).toTimeString().split(' ')[0].slice(0, -3);
    }

    UpdateCurrentProgress() {
        let progress = document.getElementById('progress-control');
        progress.value = this.player.currentTime;

        let currentTime = document.getElementById('current-time');
        let duration = document.getElementById('duration');
        currentTime.innerHTML = this.#FormatTime(this.player.currentTime);
        duration.innerHTML = this.#FormatTime(this.player.duration);
    }

    #FormatTime(time, includeSeconds = true) {
        let hours = Math.floor(time / 3600);
        let minutes = Math.floor(time % 3600 / 60);
        if (!includeSeconds) {
            return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}`;
        }

        let seconds = Math.floor(time % 3600 % 60);
        return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
    }

    VolumeUp() {
        this.player.volume = Math.min(this.player.volume + this.volumeIncrement, 1);
        this.UpdateVolumeControl();
    }

    VolumeDown() {
        this.player.volume = Math.max(this.player.volume - this.volumeIncrement, 0);
        this.UpdateVolumeControl();
    }

    Mute() {
        this.player.muted = !this.player.muted;
        this.UpdateVolumeControl();
    }

    UpdateVolumeControl() {
        let volumeControl = document.getElementById('volume-control');
        volumeControl.value = this.player.volume;

        let volumeIcon = document.getElementById('volume-icon');
        if (this.player.muted) {
            volumeIcon.innerHTML = "volume_off";
        } else if (this.player.volume < 0.5) {
            volumeIcon.innerHTML = "volume_down";
        } else {
            volumeIcon.innerHTML = "volume_up";
        }
    }


    ToggleFullScreen() {
        if (document.fullscreenElement) {
            if (document.exitFullscreen) {
                document.exitFullscreen();
            }
            return;
        }

        if (document.documentElement.requestFullscreen) {
            document.documentElement.requestFullscreen();
        }

    }



}