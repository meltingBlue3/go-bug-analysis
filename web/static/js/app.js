(function () {
    'use strict';

    // ===========================
    // DOM References
    // ===========================
    var uploadSection = document.getElementById('upload-section');
    var dashboard = document.getElementById('dashboard');
    var csvFileInput = document.getElementById('csv-file');
    var uploadStatus = document.getElementById('upload-status');

    // ===========================
    // Status Display
    // ===========================
    function showStatus(message, type) {
        if (!uploadStatus) return;
        type = type || 'info';
        uploadStatus.innerHTML = '<div class="status-' + type + '">' + escapeHtml(message) + '</div>';
    }

    function clearStatus() {
        if (uploadStatus) {
            uploadStatus.innerHTML = '';
        }
    }

    // ===========================
    // HTML Escape Utility
    // ===========================
    function escapeHtml(text) {
        var div = document.createElement('div');
        div.appendChild(document.createTextNode(text));
        return div.innerHTML;
    }

    // ===========================
    // File Size Formatting
    // ===========================
    function formatFileSize(bytes) {
        if (bytes === 0) return '0 B';
        var units = ['B', 'KB', 'MB', 'GB'];
        var i = Math.floor(Math.log(bytes) / Math.log(1024));
        return (bytes / Math.pow(1024, i)).toFixed(1) + ' ' + units[i];
    }

    // ===========================
    // Dashboard Toggle
    // ===========================
    function showDashboard() {
        if (uploadSection) uploadSection.style.display = 'none';
        if (dashboard) dashboard.style.display = 'block';
    }

    function showUpload() {
        if (uploadSection) uploadSection.style.display = '';
        if (dashboard) dashboard.style.display = 'none';
    }

    // ===========================
    // CSV Upload Handler
    // ===========================
    function uploadCSV(file) {
        showStatus('文件已选择: ' + file.name + ' (' + formatFileSize(file.size) + ')，上传功能将在下一步实现', 'info');
    }

    // ===========================
    // File Input Listener
    // ===========================
    function initFileInput() {
        if (!csvFileInput) return;

        csvFileInput.addEventListener('change', function (e) {
            var file = e.target.files && e.target.files[0];
            if (!file) {
                clearStatus();
                return;
            }

            // Check file extension
            var name = file.name.toLowerCase();
            if (name.indexOf('.csv') !== name.length - 4) {
                showStatus('请选择 .csv 格式的文件', 'error');
                csvFileInput.value = '';
                return;
            }

            uploadCSV(file);
        });
    }

    // ===========================
    // Initialization
    // ===========================
    function init() {
        initFileInput();
    }

    // Wait for DOM ready
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }

    // ===========================
    // Public API
    // ===========================
    window.BugAnalysis = {
        showStatus: showStatus,
        clearStatus: clearStatus,
        showDashboard: showDashboard,
        showUpload: showUpload
    };

})();
