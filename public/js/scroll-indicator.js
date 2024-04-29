export class ScrollIndicator {
    scrollProgress;
    content;

    constructor() {
        this.scrollProgress = document.getElementsByClassName('scrollProgress')[0];
        this.content = document.getElementsByClassName('content')[0];
        if (this.scrollProgress.length === 0) {
            return;
        }

        document.addEventListener('scroll', this.RefreshProgress.bind(this));
    }

    RefreshProgress() {
        const pageHeight = document.documentElement.scrollHeight - document.documentElement.clientHeight;
        const scrollHeight = window.scrollY;
        const scrolledPercent = scrollHeight / pageHeight;
        this.scrollProgress.style.width = `${scrolledPercent * 100}%`;
        if (scrolledPercent * 100 < 0.1) {
            this.scrollProgress.querySelector('.scrollProgressThumb').style.display = "none";
        } else {
            this.scrollProgress.querySelector('.scrollProgressThumb').style.display = "block";
        }
    }
}