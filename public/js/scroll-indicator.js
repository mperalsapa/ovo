export class ScrollIndicator {
    scrollProgress;
    content;

    constructor() {
        let scrollProgress = document.getElementsByClassName('scrollProgress');
        if (scrollProgress.length === 0) {
            return;
        }
        this.scrollProgress = scrollProgress[0];

        let content = document.getElementsByClassName('content');
        if (content.length === 0) {
            return;
        }

        this.content = content[0];
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