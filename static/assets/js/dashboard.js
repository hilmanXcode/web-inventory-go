// 1. Toggle Sidebar (Sangat berguna untuk tampilan HP)
const sidebar = document.getElementById('sidebar');
const toggleBtn = document.getElementById('toggleBtn');
const mainWrapper = document.querySelector('.main-wrapper');
const closeBtnMobile = document.querySelector('.close-btn-mobile')

toggleBtn.addEventListener('click', () => {
// Jika di layar HP, kita tambahkan class 'active' untuk memunculkan
if(window.innerWidth <= 768) {
    sidebar.classList.toggle('active');
} else {
    // Di Desktop, kita sembunyikan dengan mengurangi margin/width
    if (sidebar.style.width === '0px') {
        sidebar.style.width = '260px';
    } else {
        sidebar.style.width = '0px';
        sidebar.style.overflow = 'hidden';
    }
}
});

closeBtnMobile.addEventListener('click', () => {
    sidebar.classList.toggle('active')
})

// 2. State Menu Active (Agar menu yang diklik menyala biru)
const menuLinks = document.querySelectorAll('.menu-link');

menuLinks.forEach(link => {
link.addEventListener('click', function() {
    // Hapus class active dari semua menu
    menuLinks.forEach(item => item.classList.remove('active'));
    // Tambahkan class active ke menu yang sedang diklik
    this.classList.add('active');
});
});

// 3. Efek Logout sederhana
// document.getElementById('btnLogout').addEventListener('click', function(e) {
// e.preventDefault();
// alert('Anda telah berhasil Logout!');
// // Di implementasi asli, arahkan ke halaman login:
// // window.location.href = '/login'; 
// });

// Fungsi untuk membuka modal
function openModal(modalId) {
    const modal = document.getElementById(modalId);
    if(modal) {
        modal.style.display = 'flex';
        // Sedikit delay agar transisi CSS berjalan
        setTimeout(() => modal.classList.add('show'), 10);
    }
}

// Fungsi untuk menutup modal
function closeModal(modalId) {
    const modal = document.getElementById(modalId);
    if(modal) {
        modal.classList.remove('show');
        // Tunggu animasi CSS selesai baru hilangkan elemen
        setTimeout(() => modal.style.display = 'none', 300);
    }
}

// Menutup modal jika user klik area gelap di luar form
window.addEventListener('click', function(e) {
    if (e.target.classList.contains('modal-overlay')) {
        closeModal(e.target.id);
    }
});