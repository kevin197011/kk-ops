document.addEventListener('DOMContentLoaded', function() {
    // 验证密码确认
    const passwordField = document.getElementById('password');
    const confirmPasswordField = document.getElementById('confirm_password');
    const userForm = document.getElementById('userForm');

    if (userForm) {
        userForm.addEventListener('submit', function(e) {
            if (passwordField.value !== confirmPasswordField.value) {
                e.preventDefault();
                alert('两次输入的密码不一致，请重新输入');
                confirmPasswordField.value = '';
                confirmPasswordField.focus();
            }

            // 验证密码长度
            if (passwordField.value.length < 8) {
                e.preventDefault();
                alert('密码长度必须至少为8个字符');
                passwordField.focus();
            }
        });
    }
});