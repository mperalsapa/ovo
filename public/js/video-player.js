import { PlayerIframe } from "./player-iframe.js";
import { Routes } from "./routes.js";

export class VideoPlayer {
    volumeIncrement = 0.05;     // 5%
    videoScrubIncrement = 5;    // seconds

    syncConnection;

    syncLatency;
    syncLatencySended;
    syncLatencyInterval;
    syncMessageDebouncerTimer; // seconds

    lastBufferState = Date.now();
    isBuffering = false;
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


        console.log("Video Player loaded")

        if (this.player.readyState < 3) {
            this.player.addEventListener('loadedmetadata', () => {
                console.log("Metadata Loadeded")
                this.#Init()
            });
        } else {
            this.#Init()
        }
    }

    #Init() {
        // Check if this video is hooked to a syncplay group
        if (this.player.dataset.enabledsyncplay == "true" && this.syncConnection === undefined) {
            this.#InitSyncplay();
            return
        }

        this.#UpdateUI();

        this.RequestCanPlay();

        // this.lastBufferingTime = Date.now();
        this.StartBufferCheckerInterval();
    }

    StartBufferCheckerInterval() {
        console.log("Starting buffer checker")
        this.bufferCheckerInterval = setInterval(() => {

            let availBuffer = this.GetAvailableBuffer();
            if (availBuffer < 1 && !this.isBuffering) {
                this.lastBufferState = Date.now();
                this.Buffering();
            }

            if (availBuffer > 1 && this.isBuffering && Date.now() - this.lastBufferState > 1000) {
                this.lastBufferState = Date.now();
                this.Canplay();
            }

        }, 100);
    }

    async #InitSyncplay() {
        this.syncConnection = new WebSocket(Routes.Routes.Websocket);
        // Hook to syncplay events
        this.syncConnection.onmessage = this.#SyncPlayListener.bind(this);
        await this.#AwaitSyncplayConnection(this.syncConnection);
        console.log("Connected to syncplay", this.syncConnection.readyState)
        this.#StartLatencyMeasurement();
        this.#Init()
    }

    #AwaitSyncplayConnection(websocket) {
        return new Promise((resolve, reject) => {
            const maxNumberOfAttempts = 10;
            const intervalTime = 2000; //ms

            let currentAttempt = 0;
            let interval = setInterval(() => {
                currentAttempt++;
                if (currentAttempt > maxNumberOfAttempts) {
                    clearInterval(interval);
                    reject(new Error("Max number of attempts reached"));
                } else if (websocket.readyState === websocket.OPEN) {
                    clearInterval(interval);
                    resolve();
                }
            }, intervalTime);
        });
    }

    #StartLatencyMeasurement() {
        // this.syncLatencyInterval = setInterval(() => {
        //     this.syncLatencySended = Date.now();
        //     this.#SendWebsocketMessage({
        //         event: "ping",
        //     });
        // }, 5000);
    }

    #SendWebsocketMessage(message) {
        // clearTimeout(this.syncMessageDebouncerTimer);
        // this.syncMessageDebouncerTimer = setTimeout(() => {
        // }, 500);
        console.log("Sending to server: ", message)
        this.syncConnection.send(JSON.stringify(message));
    }

    #SyncPlayListener(event) {
        let data = JSON.parse(event.data);
        console.log("Received from server: " + event.data);
        switch (data.event) {
            case "pong":
                this.syncLatency = Date.now() - this.syncLatencySended;
                console.log("Latency: ", this.syncLatency)
                break;
            case "play":
                console.log("Playing from: ", data.StartedFrom)
                this.player.currentTime = this.GetCurrentTime(data.StartedFrom, data.StartedAt);
                this.#Play();
                break;
            case "pause":
                console.log("Paused")
                this.player.currentTime = data.StartedFrom;
                this.player.pause();
                break;
            case "seek":
                console.log("Seek to new starting point: ", data.StartedFrom)
                this.player.currentTime = this.GetCurrentTime(data.StartedFrom, data.StartedAt);
                break;
            case "newItem":
                console.log("New item: ", data.Item.Title, "(ID: ", data.Item.ID, ")");
                // load current iframe into new url containing ID
                window.location.href = Routes.Routes.Player;
                break;
            case "buffering":
                console.log("Other player is buffering, waiting there: ", data.StartedFrom)
                this.player.pause();
                this.player.currentTime = data.StartedFrom;
                break;
            case "canplay":
                console.log("Other player can play, waiting there: ", data.StartedFrom)
                this.player.currentTime = this.GetCurrentTime(data.StartedFrom, data.StartedAt);
                break;
            default:
                console.log("Event not handled: ", data.event);
                break;
        }
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
        this.player.addEventListener('click', this.RequestPlay.bind(this));
        this.player.addEventListener('dblclick', this.ToggleFullScreen.bind(this));
        this.buttons.play.addEventListener('click', this.RequestPlay.bind(this));
        this.buttons.rewind.addEventListener('click', this.RequestRewind.bind(this));
        this.buttons.forward.addEventListener('click', this.RequestForward.bind(this));
        this.buttons.mute.addEventListener('click', this.Mute.bind(this));
        this.buttons.fullScreen.addEventListener('click', this.ToggleFullScreen.bind(this));

        // keyboard listeners
        document.addEventListener('keydown', this.#KeyboardListener.bind(this));

        // video listeners
        //   - play
        this.player.addEventListener('play', function () {
            this.UpdatePlayButton("pause")
            this.UpdateEndsAt();
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
            this.SaveVolume();
            this.UpdateVolumeControl();
        }.bind(this));

        //   - progress
        this.player.addEventListener('timeupdate', this.UpdateCurrentProgress.bind(this));
        let progressControl = document.getElementById('progress-control');

        progressControl.addEventListener('input', function (e) {
            let seekedTime = e.target.value / progressControl.max * this.player.duration;
            this.RequestSeek(seekedTime);
        }.bind(this));

        //  - Waiting & CanPlay
        // this.player.addEventListener('waiting', this.Buffering.bind(this));
        // this.player.addEventListener('canplay', this.Canplay.bind(this));
    }

    #KeyboardListener(e) {
        switch (e.code) {
            case 'Space':
                this.RequestPlay();
                break;
            case 'ArrowLeft':
                this.RequestRewind();
                break;
            case 'ArrowRight':
                this.RequestForward();
                break;
            case 'Home':
                this.RequestSeek(0);
                break;
            case 'End':
                this.RequestSeek(this.player.duration);
                break;
            case 'ArrowUp':
                this.VolumeUp();
                break;
            case 'ArrowDown':
                this.VolumeDown();
                break;
            case 'Digit1':
                this.RequestSeek(this.player.duration * 0.1);
                break;
            case 'Digit2':
                this.RequestSeek(this.player.duration * 0.2);
                break;
            case 'Digit3':
                this.RequestSeek(this.player.duration * 0.3);
                break;
            case 'Digit4':
                this.RequestSeek(this.player.duration * 0.4);
                break;
            case 'Digit5':
                this.RequestSeek(this.player.duration * 0.5);
                break;
            case 'Digit6':
                this.RequestSeek(this.player.duration * 0.6);
                break;
            case 'Digit7':
                this.RequestSeek(this.player.duration * 0.7);
                break;
            case 'Digit8':
                this.RequestSeek(this.player.duration * 0.8);
                break;
            case 'Digit9':
                this.RequestSeek(this.player.duration * 0.9);
                break;
            case 'KeyF':
                this.ToggleFullScreen();
                break;
            case 'Period':
                this.RequestPause();
                this.RequestSeek(this.player.currentTime + this.GetFrameLength());
                console.log("Frame Time: ", this.GetFrameLength())
                break;
            case 'Comma':
                this.RequestPause();
                this.RequestSeek(this.player.currentTime - this.GetFrameLength());
                break;
            default:
                console.log(e.code + " not handled")
                break;
        }
    }

    #UpdateUI() {
        this.UpdatePlayButton();
        this.UpdateEndsAt();
        this.UpdateCurrentProgress();
        this.UpdateVolumeControl();
    }

    Buffering() {
        let actualDate = new Date();
        this.isBuffering = true;
        console.log("Available buffer: ", this.GetAvailableBuffer(), " seconds")
        console.log("Buffering... " + actualDate.getHours() + ":" + actualDate.getMinutes() + ":" + actualDate.getSeconds() + "." + actualDate.getMilliseconds())
        this.#SendWebsocketMessage({
            event: "buffering",
            StartedFrom: this.player.currentTime
        })
    }

    GetAvailableBuffer() {
        let currentTime = this.player.currentTime;
        let buffered = this.player.buffered;
        let bufferEnd = currentTime;
        for (let i = 0; i < buffered.length; i++) {
            if (buffered.start(i) <= currentTime && buffered.end(i) >= currentTime) {
                // console.log("Buffered: ", buffered.start(i), buffered.end(i))
                if (buffered.end(i) > bufferEnd) {
                    bufferEnd = buffered.end(i);
                }
                break;
            }
        }
        return bufferEnd - currentTime;
    }

    Canplay() {
        this.isBuffering = false;
        let actualDate = new Date();
        console.log("Can play :D " + actualDate.getHours() + ":" + actualDate.getMinutes() + ":" + actualDate.getSeconds() + "." + actualDate.getMilliseconds())
        this.#SendWebsocketMessage({
            event: "canplay",
            StartedFrom: this.player.currentTime
        })
    }

    // This function requires an number (startedFrom) and a Unix millisecond timestamp (startedAt)
    // Returns the elapsed time + the offset of the startedFrom to the current time
    GetCurrentTime(StartedFrom, StartedAt) {
        return StartedFrom + (Date.now() - StartedAt) / 1000;
    }

    RequestCanPlay() {
        if (this.syncConnection) {
            this.#SendWebsocketMessage({
                event: "requestPlay",
            });
        } else {
            this.#TogglePlayPause();
        }
    }

    RequestPlay() {
        if (this.syncConnection) {
            let newAction = this.player.paused ? "play" : "pause";
            this.#SendWebsocketMessage({
                event: newAction,
                StartedFrom: this.player.currentTime
            });
        } else {
            this.#TogglePlayPause();
        }
    }

    RequestPause() {
        if (this.syncConnection) {
            this.#SendWebsocketMessage({
                event: "pause",
            });
        } else {
            this.#TogglePlayPause();
        }
    }

    #TogglePlayPause() {
        this.player.paused ? this.#Play() : this.player.pause();
    }

    #Play() {
        console.log("Duration: ", this.player.duration, "Current Time: ", this.player.currentTime)
        if (this.player.currentTime === this.player.duration) {
            // This is required because some browsers (firefox) wont start from the beginning if the video is already ended
            this.RequestSeek(0);
            return
        }

        this.player.play();
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

    RequestSeek(runtime) {
        if (this.syncConnection) {
            this.#SendWebsocketMessage({
                event: "seek",
                StartedFrom: runtime
            });
        } else {
            this.#Seek(runtime);
        }
    }

    #Seek(runtime) {
        console.log("Seeking to: ", runtime)
        this.player.currentTime = runtime;
    }

    RequestRewind() {
        if (this.syncConnection) {
            this.RequestSeek(this.player.currentTime - this.videoScrubIncrement);
        } else {
            this.#Rewind();
        }
    }

    #Rewind() {
        this.#Seek(this.player.currentTime - this.videoScrubIncrement);
    }

    RequestForward() {
        if (this.syncConnection) {
            this.RequestSeek(this.player.currentTime + this.videoScrubIncrement);
        } else {
            this.#Forward();
        }
    }


    #Forward() {
        this.#Seek(this.player.currentTime + this.videoScrubIncrement);
    }


    UpdateEndsAt() {
        let endsAt = document.getElementById('ends-at');
        let currentTime = Date.now();
        currentTime += (this.player.duration - this.player.currentTime) * 1000;

        endsAt.innerHTML = new Date(currentTime).toTimeString().split(' ')[0].slice(0, -3);
    }

    UpdateCurrentProgress() {
        let progress = document.getElementById('progress-control');
        progress.value = (this.player.currentTime / this.player.duration) * progress.max;

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
        this.SaveVolume();
        this.UpdateVolumeControl();
    }

    VolumeDown() {
        this.player.volume = Math.max(this.player.volume - this.volumeIncrement, 0);
        this.SaveVolume();
        this.UpdateVolumeControl();
    }

    Mute() {
        this.player.muted = !this.player.muted;
        this.UpdateVolumeControl();
    }

    UpdateVolumeControl() {
        // Load volume from localstorage
        this.player.volume = Math.min(Math.max(this.LoadVolume(), 0), 1);
        console.log("Volume: ", this.player.volume)
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

    SaveVolume() {
        localStorage.setItem('volume', this.player.volume);
    }

    LoadVolume() {
        let volume = localStorage.getItem('volume');
        if (volume) {
            return parseFloat(volume);
        }
        return 1;
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

    GetFrameLength() {
        // we asume videos are 30fps. While is wide known that movies are filmed in 24fps, some
        // are filmed in 30fps, and is better to get smaller frame time than bigger (30fps -> 33ms vs 24fps -> 41ms)
        return 1 / 30;
    }


}