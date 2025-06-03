// 全局变量
let authKey = '';

// 页面加载完成后检查认证状态
document.addEventListener('DOMContentLoaded', function() {
    // 检查是否已有保存的密钥
    const savedKey = sessionStorage.getItem('authKey');
    if (savedKey) {
        authKey = savedKey;
        checkAuthAndProceed();
    } else {
        showAuthModal();
    }
    
    document.getElementById('authKeyInput').addEventListener('keypress', function(e) {
        if (e.key === 'Enter') {
            authenticate();
        }
    });
});

// 显示验证模态框
function showAuthModal() {
    document.getElementById('authModal').style.display = 'flex';
    document.getElementById('mainContent').classList.remove('authenticated');
    document.getElementById('authKeyInput').focus();
}

// 验证密钥
function authenticate() {
    const inputKey = document.getElementById('authKeyInput').value.trim();
    const errorElement = document.getElementById('authError');
    const btnElement = document.getElementById('authBtn');
    
    if (!inputKey) {
        showAuthError('请输入访问密钥');
        return;
    }
    
    // 禁用按钮，显示加载状态
    btnElement.disabled = true;
    btnElement.textContent = '验证中...';
    errorElement.textContent = '';
    
    // 测试密钥是否正确（发送一个测试请求）
    fetch('/ping', {
        method: 'GET',
        headers: {
            'X-Auth-Key': inputKey
        }
    })
    .then(response => {
        if (response.ok) {
            // 密钥正确
            authKey = inputKey;
            sessionStorage.setItem('authKey', authKey);
            hideAuthModal();
            loadMainContent();
        } else if (response.status === 401) {
            // 密钥错误
            showAuthError('访问密钥错误，请重新输入');
        } else {
            showAuthError('服务器连接失败，请稍后重试');
        }
    })
    .catch(error => {
        console.error('验证失败:', error);
        showAuthError('网络错误，请检查连接');
    })
    .finally(() => {
        btnElement.disabled = false;
        btnElement.textContent = '验证并进入';
    });
}

function checkAuthAndProceed() {
    fetch('/ping', {
        method: 'GET',
        headers: {
            'X-Auth-Key': authKey
        }
    })
    .then(response => {
        if (response.ok) {
            hideAuthModal();
            loadMainContent();
        } else {
            sessionStorage.removeItem('authKey');
            authKey = '';
            showAuthModal();
        }
    })
    .catch(error => {
        console.error('认证检查失败:', error);
        showAuthModal();
    });
}

function showAuthError(message) {
    const errorElement = document.getElementById('authError');
    errorElement.textContent = message;
    
    document.getElementById('authKeyInput').value = '';
    document.getElementById('authKeyInput').focus();
    
    setTimeout(() => {
        errorElement.textContent = '';
    }, 3000);
}

function hideAuthModal() {
    document.getElementById('authModal').style.display = 'none';
    document.getElementById('mainContent').classList.add('authenticated');
}

function loadMainContent() {
    document.getElementById('result').style.display = 'none';
    document.getElementById('urlResult').style.display = 'none';
    
    loadImageList();
}

// 退出登录
function logout() {
    if (confirm('确定要退出登录吗？')) {
        sessionStorage.removeItem('authKey');
        authKey = '';
        showAuthModal();
    }
}

function authenticatedFetch(url, options = {}) {
    const headers = {
        'X-Auth-Key': authKey,
        ...options.headers
    };
    
    return fetch(url, {
        ...options,
        headers
    }).then(response => {
        // 如果返回401，说明密钥失效
        if (response.status === 401) {
            sessionStorage.removeItem('authKey');
            authKey = '';
            showAuthModal();
            throw new Error('认证失败');
        }
        return response;
    });
}

const uploadArea = document.querySelector('.upload-area');
const fileInput = document.getElementById('fileInput');
const progressBar = document.querySelector('.progress-bar');
const progressFill = document.querySelector('.progress-fill');
const result = document.getElementById('result');
const urlResult = document.getElementById('urlResult');
const imageUrl = document.getElementById('imageUrl');

// 拖拽上传
uploadArea.addEventListener('dragover', (e) => {
    e.preventDefault();
    uploadArea.classList.add('dragover');
});

uploadArea.addEventListener('dragleave', () => {
    uploadArea.classList.remove('dragover');
});

uploadArea.addEventListener('drop', (e) => {
    e.preventDefault();
    uploadArea.classList.remove('dragover');
    const files = e.dataTransfer.files;
    if (files.length > 0) {
        uploadFile(files[0]);
    }
});

fileInput.addEventListener('change', (e) => {
    if (e.target.files.length > 0) {
        uploadFile(e.target.files[0]);
    }
});

function uploadFile(file, retryCount = 0) {
    const maxRetries = 2;
    const formData = new FormData();
    formData.append('image', file);

    progressBar.style.display = 'block';
    progressFill.style.width = '0%';
    result.style.display = 'none';
    urlResult.style.display = 'none';

    // 显示上传状态
    if (retryCount > 0) {
        showResult({success: false, message: `正在重试上传... (${retryCount}/${maxRetries})`});
    }

    // 模拟进度条
    let progress = 0;
    const progressInterval = setInterval(() => {
        progress += Math.random() * 15; 
        if (progress > 85) progress = 85; 
        progressFill.style.width = progress + '%';
    }, 200); 

    const controller = new AbortController();
    const timeoutId = setTimeout(() => {
        controller.abort();
    }, 300000); 

    authenticatedFetch('/upload', {
        method: 'POST',
        body: formData,
        signal: controller.signal
    })
    .then(response => {
        clearTimeout(timeoutId);
        return response.json();
    })
    .then(data => {
        clearInterval(progressInterval);
        progressFill.style.width = '100%';

        setTimeout(() => {
            progressBar.style.display = 'none';
            showResult(data);
            if (data.success) {
                loadImageList();
            }
        }, 500);
    })
    .catch(error => {
        clearTimeout(timeoutId);
        clearInterval(progressInterval);
        progressBar.style.display = 'none';

        if (error.message !== '认证失败') {
            if (retryCount < maxRetries && (
                error.name === 'AbortError' ||
                error.message.includes('timeout') ||
                error.message.includes('network') ||
                error.message.includes('fetch')
            )) {
                console.log(`上传失败，准备重试 (${retryCount + 1}/${maxRetries}):`, error.message);
                setTimeout(() => {
                    uploadFile(file, retryCount + 1);
                }, 2000); // 2秒后重试
            } else {
                let errorMessage = '上传失败';
                if (error.name === 'AbortError') {
                    errorMessage = '上传超时，请检查网络连接或尝试上传更小的文件';
                } else if (error.message.includes('timeout')) {
                    errorMessage = '网络超时，请检查网络连接';
                } else if (error.message.includes('Failed to fetch')) {
                    errorMessage = '网络连接失败，请检查网络';
                } else {
                    errorMessage = '上传失败: ' + error.message;
                }
                showResult({success: false, message: errorMessage});
            }
        }
    });
}

function showResult(data) {
    result.style.display = 'block';
    result.className = 'result ' + (data.success ? 'success' : 'error');
    result.textContent = data.message;
    
    if (data.success && data.url) {
        imageUrl.value = data.url;
        urlResult.style.display = 'block';
    }
}

function copyUrl(event) {
    navigator.clipboard.writeText(imageUrl.value).then(() => {
        const btn = event.target;
        const originalText = btn.textContent;
        btn.textContent = '已复制!';
        setTimeout(() => {
            btn.textContent = originalText;
        }, 2000);
    }).catch(() => {
        imageUrl.select();
        document.execCommand('copy');
        const btn = event.target;
        const originalText = btn.textContent;
        btn.textContent = '已复制!';
        setTimeout(() => {
            btn.textContent = originalText;
        }, 2000);
    });
}

function openImage() {
    window.open(imageUrl.value, '_blank');
}

function loadImageList() {
    document.getElementById('loading').style.display = 'block';
    document.getElementById('imageGrid').innerHTML = '';
    
    authenticatedFetch('/list')
    .then(response => response.json())
    .then(images => {
        document.getElementById('loading').style.display = 'none';
        const grid = document.getElementById('imageGrid');
        
        images.slice(0, 20).forEach(image => {
            const item = document.createElement('div');
            item.className = 'image-item';
            item.innerHTML = 
                '<img src="' + image.url + '" alt="' + image.filename + '" loading="lazy">' +
                '<div class="image-info">' +
                    '<div><strong>' + image.filename + '</strong></div>' +
                    '<div>大小: ' + formatFileSize(image.size) + '</div>' +
                    '<div>时间: ' + image.date + '</div>' +
                    '<div class="image-url" onclick="copyImageUrl(\'' + window.location.origin + image.url + '\', event)" title="点击复制链接">' +
                        window.location.origin + image.url +
                    '</div>' +
                '</div>';
            grid.appendChild(item);
        });
    })
    .catch(error => {
        document.getElementById('loading').style.display = 'none';
        if (error.message !== '认证失败') {
            console.error('加载图片列表失败:', error);
        }
    });
}

function copyImageUrl(url, event) {
    navigator.clipboard.writeText(url).then(() => {
        event.target.style.background = '#d4edda';
        setTimeout(() => {
            event.target.style.background = '#1a1a1a';
        }, 1000);
    }).catch(() => {
        console.log('复制失败，URL:', url);
    });
}

function formatFileSize(bytes) {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
}
