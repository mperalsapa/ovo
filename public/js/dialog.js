export class Dialog {
    dialog;
    parent;
    constructor(parent) {
        if (!parent) {
            throw new Error('Dialog: Parent is required');
        }
        this.parent = parent;

        // Creating dialog
        this.dialog = document.createElement('dialog');
        this.dialog.classList.add('dialog');
        document.documentElement.appendChild(this.dialog);

        // Trigger open
        this.parent.addEventListener('click', (e) => { this.OpenDialog(e) });

        // Trigger close
        this.dialog.addEventListener('click', (e) => {
            var rect = this.dialog.getBoundingClientRect();
            var isInDialog = (rect.top <= e.clientY && e.clientY <= rect.top + rect.height &&
                rect.left <= e.clientX && e.clientX <= rect.left + rect.width);
            if (!isInDialog) {
                this.dialog.close();
            }
        });
    }

    OpenDialog(e) {
        let x = e.x;
        let y = e.y;
        let maxX = this.dialog.clientWidth;
        let maxY = this.dialog.clientHeight;

        this.dialog.style.marginLeft = `${x - maxX / 2}px`;
        this.dialog.style.marginTop = `${y - maxY / 2}px`;
        this.dialog.showModal();
    }

    SetContent(content) {
        this.dialog.innerHTML = content;
    }
}

export class DialogList extends Dialog {
    container;
    header;
    content;
    constructor(parent) {
        super(parent)

        this.container = document.createElement('div');
        this.container.classList.add('dialog-list');
        this.header = document.createElement('h2');
        this.content = document.createElement('div');

        this.dialog.appendChild(this.container);
        this.container.appendChild(this.header);
        this.container.appendChild(this.content);
    }

    SetHeader(header) {
        this.header.innerText = header;
    }

    SetList(elementList) {
        if (Array.isArray(elementList)) {
            elementList.forEach(element => {
                this.content.appendChild(element);
            });
        }
    }
}