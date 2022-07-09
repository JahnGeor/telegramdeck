const signInBtn = document.querySelector(".signin-btn");
const signUpBtn = document.querySelector(".signup-btn");
const formBox = document.querySelector(".form-box");
const header = document.querySelector(".header");
const body = document.body;

signUpBtn.addEventListener('click', function () {
    formBox.classList.add('active');
    body.classList.add('active');
    header.classList.add('active');
})

signInBtn.addEventListener('click', function () {
    formBox.classList.remove('active');
    body.classList.remove('active');
    header.classList.remove('active');
})