// 全局变量
let currentPage = 1;
let token = localStorage.getItem('token');

// 更新导航栏状态
function updateNavbar() {
    const navbarNav = document.querySelector('.navbar-nav:last-child');
    if (token) {
        navbarNav.innerHTML = `
            <a class="nav-link" href="/admin">管理后台</a>
            <a class="nav-link" href="#" onclick="logout()">退出</a>
        `;
    } else {
        navbarNav.innerHTML = `
            <a class="nav-link" href="/login">登录</a>
            <a class="nav-link" href="/register">注册</a>
        `;
    }
}

// 加载文章列表
async function loadPosts(page = 1) {
    try {
        const response = await fetch(`/api/posts?page=${page}`);
        const data = await response.json();
        
        const postsContainer = document.getElementById('posts-container');
        postsContainer.innerHTML = '';
        
        data.posts.forEach(post => {
            const postElement = document.createElement('div');
            postElement.className = 'card-body border-bottom';
            postElement.innerHTML = `
                <h2 class="card-title">
                    <a href="/posts/${post.ID}" class="text-decoration-none">${post.title}</a>
                </h2>
                <p class="card-text text-muted">
                    <small>作者: ${post.user.username} | 
                    分类: ${post.category.name} | 
                    发布时间: ${new Date(post.publish_at).toLocaleDateString()}</small>
                </p>
                <p class="card-text">${post.content.substring(0, 200)}...</p>
                <a href="/posts/${post.ID}" class="btn btn-primary btn-sm">阅读更多</a>
            `;
            postsContainer.appendChild(postElement);
        });
        
        // 更新分页
        updatePagination(data.total);
    } catch (error) {
        console.error('加载文章失败:', error);
    }
}

// 更新分页控件
function updatePagination(total) {
    const totalPages = Math.ceil(total / 10);
    const pagination = document.querySelector('.pagination');
    pagination.innerHTML = '';
    
    // 上一页
    pagination.innerHTML += `
        <li class="page-item ${currentPage === 1 ? 'disabled' : ''}">
            <a class="page-link" href="#" onclick="loadPosts(${currentPage - 1})">上一页</a>
        </li>
    `;
    
    // 页码
    for (let i = 1; i <= totalPages; i++) {
        pagination.innerHTML += `
            <li class="page-item ${currentPage === i ? 'active' : ''}">
                <a class="page-link" href="#" onclick="loadPosts(${i})">${i}</a>
            </li>
        `;
    }
    
    // 下一页
    pagination.innerHTML += `
        <li class="page-item ${currentPage === totalPages ? 'disabled' : ''}">
            <a class="page-link" href="#" onclick="loadPosts(${currentPage + 1})">下一页</a>
        </li>
    `;
}

// 加载分类列表
async function loadCategories() {
    try {
        const response = await fetch('/api/categories');
        const categories = await response.json();
        
        const categoriesList = document.getElementById('categories-list');
        categoriesList.innerHTML = categories.map(category => `
            <li>
                <a href="/categories/${category.id}" class="text-decoration-none">
                    ${category.name} (${category.posts.length})
                </a>
            </li>
        `).join('');
    } catch (error) {
        console.error('加载分类失败:', error);
    }
}

// 退出登录
function logout() {
    localStorage.removeItem('token');
    token = null;
    updateNavbar();
    window.location.href = '/';
}

// 页面加载完成后执行
document.addEventListener('DOMContentLoaded', () => {
    updateNavbar();
    loadPosts();
    loadCategories();
});