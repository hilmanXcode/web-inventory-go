const container = document.getElementById('container');

// Tombol versi Desktop
const btnGoRegister = document.getElementById('btn-go-register');
const btnGoLogin = document.getElementById('btn-go-login');

// Tombol versi Mobile
const mobileGoRegister = document.getElementById('mobile-go-register');
const mobileGoLogin = document.getElementById('mobile-go-login');

// Logic animasi: menambah atau menghapus class 'sign-up-mode' pada container
btnGoRegister.addEventListener('click', () => {
    container.classList.add('sign-up-mode');
});

btnGoLogin.addEventListener('click', () => {
    container.classList.remove('sign-up-mode');
});

// Logic untuk versi mobile
mobileGoRegister.addEventListener('click', (e) => {
    e.preventDefault(); // Mencegah reload halaman
    container.classList.add('sign-up-mode');
});

mobileGoLogin.addEventListener('click', (e) => {
    e.preventDefault();
    container.classList.remove('sign-up-mode');
});

// Mencegah form reload saat tombol di klik (hanya untuk demo)
// document.getElementById('form1').addEventListener('submit', (e) => e.preventDefault());
// document.getElementById('form2').addEventListener('submit', (e) => e.preventDefault());