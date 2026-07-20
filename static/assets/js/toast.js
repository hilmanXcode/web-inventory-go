class Toast {
    constructor() {
        // Buat container jika belum ada di DOM
        this.container = document.getElementById('toast-container');
        if (!this.container) {
            this.container = document.createElement('div');
            this.container.id = 'toast-container';
            document.body.appendChild(this.container);
        }

        // Kumpulan icon SVG sederhana agar tidak perlu library icon eksternal
        this.icons = {
            success: '✔️', // Anda bisa ganti dengan kode SVG
            error: '✖️',
            info: 'ℹ️',
            warning: '⚠️'
        };
    }

    /**
     * Method utama untuk memanggil toast
     * @param {Object} options - { title, message, type, duration }
     */
    show(options) {
        const { title = 'Notifikasi', message = '', type = 'info', duration = 3000 } = options;

        // Buat elemen utama toast
        const toastEl = document.createElement('div');
        toastEl.className = `toast-item ${type}`;

        // Template HTML isi toast
        toastEl.innerHTML = `
            <div class="toast-icon">${this.icons[type]}</div>
            <div class="toast-body">
                <div class="toast-title">${title}</div>
                <div class="toast-message">${message}</div>
            </div>
            <button class="toast-close">&times;</button>
            <div class="toast-progress" style="animation-duration: ${duration}ms;"></div>
        `;

        this.container.appendChild(toastEl);

        // Beri sedikit jeda sebelum menambah class 'show' agar animasi CSS berjalan
        requestAnimationFrame(() => {
            toastEl.classList.add('show');
        });

        // Logika untuk menghapus toast
        const removeToast = () => {
            toastEl.classList.remove('show');
            toastEl.classList.add('hide');
            
            // Tunggu animasi selesai baru hapus elemen dari DOM
            toastEl.addEventListener('transitionend', () => {
                toastEl.remove();
            });
        };

        // Hapus otomatis sesuai durasi
        const timer = setTimeout(removeToast, duration);

        // Hapus manual jika tombol close ditekan
        toastEl.querySelector('.toast-close').addEventListener('click', () => {
            clearTimeout(timer);
            removeToast();
        });
    }

    // --- Shortcut Methods agar pemanggilan lebih bersih ---

    success(title, message, duration) {
        this.show({ title, message, type: 'success', duration });
    }

    error(title, message, duration) {
        this.show({ title, message, type: 'error', duration });
    }

    info(title, message, duration) {
        this.show({ title, message, type: 'info', duration });
    }

    warning(title, message, duration) {
        this.show({ title, message, type: 'warning', duration });
    }
}

// Inisialisasi instance global agar bisa langsung dipanggil
const toast = new Toast();