const API_BASE = 'http://localhost:8080';

// Общие функции
function showMessage(message, type = 'success') {
    const messageEl = document.getElementById('message');
    if (messageEl) {
        messageEl.textContent = message;
        messageEl.className = `message ${type}`;
        setTimeout(() => {
            messageEl.textContent = '';
            messageEl.className = 'message';
        }, 5000);
    }
}

function getToken() {
    return localStorage.getItem('token');
}

function setToken(token) {
    localStorage.setItem('token', token);
}

function removeToken() {
    localStorage.removeItem('token');
}

function isLoggedIn() {
    return !!getToken();
}

function redirectToLogin() {
    window.location.href = 'login.html';
}

// ===== ГЛАВНАЯ СТРАНИЦА (публичный доступ) =====
if (document.getElementById('moviesContainer') && !document.getElementById('movieForm')) {
    // Это главная страница - загружаем фильмы без авторизации

    async function loadMoviesPublic() {
        try {
            const searchParams = new URLSearchParams();

            // Добавляем параметры фильтрации
            const title = document.getElementById('searchTitle').value;
            const genre = document.getElementById('searchGenre').value;
            const director = document.getElementById('searchDirector').value;
            const sortField = document.getElementById('sortField').value;
            const sortOrder = document.getElementById('sortOrder').value;

            if (title) searchParams.append('title', title);
            if (genre) searchParams.append('genre', genre);
            if (director) searchParams.append('director', director);
            if (sortField) {
                searchParams.append('sort', sortField);
                searchParams.append('order', sortOrder);
            }

            const response = await fetch(`${API_BASE}/movies/public?${searchParams}`);
            const movies = await response.json();
            displayMoviesPublic(movies);
        } catch (error) {
            console.error('Ошибка загрузки фильмов:', error);
        }
    }

    function displayMoviesPublic(movies) {
        const container = document.getElementById('moviesContainer');

        container.innerHTML = '';

        if (movies.length === 0) {
            container.innerHTML = '<p>Фильмы не найдены</p>';
            return;
        }

        container.innerHTML = movies.map(movie => `
            <div class="movie-card">
                <h3>${movie.title}</h3>
                <p><strong>Жанр:</strong> ${movie.genre}</p>
                <p><strong>Режиссер:</strong> ${movie.director}</p>
                <p><strong>Рейтинг:</strong> <span class="rating">${movie.rating}</span>/10</p>
                <p><strong>Дата выхода:</strong> ${new Date(movie.release_date).toLocaleDateString('ru-RU')}</p>
            </div>
        `).join('');
    }

    // Применение фильтров на главной странице
    document.getElementById('applyFilters')?.addEventListener('click', loadMoviesPublic);

    // Сброс фильтров на главной странице
    document.getElementById('clearFilters')?.addEventListener('click', () => {
        document.getElementById('searchTitle').value = '';
        document.getElementById('searchGenre').value = '';
        document.getElementById('searchDirector').value = '';
        document.getElementById('sortField').value = '';
        document.getElementById('sortOrder').value = 'asc';
        loadMoviesPublic();
    });

    // Загрузка фильмов при загрузке страницы
    document.addEventListener('DOMContentLoaded', loadMoviesPublic);
}

// ===== РЕГИСТРАЦИЯ =====
if (document.getElementById('registerForm')) {
    document.getElementById('registerForm').addEventListener('submit', async (e) => {
        e.preventDefault();

        const formData = {
            username: document.getElementById('regUsername').value,
            password: document.getElementById('regPassword').value
        };

        try {
            const response = await fetch(`${API_BASE}/register`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(formData)
            });

            const data = await response.json();

            if (response.ok) {
                showMessage('Регистрация успешна! Теперь вы можете войти.', 'success');
                setTimeout(() => {
                    window.location.href = 'login.html';
                }, 2000);
            } else {
                showMessage(data.error || 'Ошибка регистрации', 'error');
            }
        } catch (error) {
            showMessage('Ошибка сети: ' + error.message, 'error');
        }
    });
}

// ===== ВХОД =====
if (document.getElementById('loginForm')) {
    document.getElementById('loginForm').addEventListener('submit', async (e) => {
        e.preventDefault();

        const formData = {
            username: document.getElementById('username').value,
            password: document.getElementById('password').value
        };

        try {
            const response = await fetch(`${API_BASE}/login`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(formData)
            });

            const data = await response.json();

            if (response.ok) {
                setToken(data.token);
                showMessage('Вход успешен!', 'success');
                setTimeout(() => {
                    window.location.href = 'movies.html';
                }, 1000);
            } else {
                showMessage(data.error || 'Ошибка входа', 'error');
            }
        } catch (error) {
            showMessage('Ошибка сети', 'error');
        }
    });
}

// ===== УПРАВЛЕНИЕ ФИЛЬМАМИ (требует авторизации) =====
if (document.getElementById('movieForm') && document.getElementById('moviesContainer')) {
    let currentEditId = null;

    // Проверка авторизации
    if (!isLoggedIn()) {
        redirectToLogin();
    }

    // Выход
    document.getElementById('logoutBtn').addEventListener('click', () => {
        removeToken();
        window.location.href = 'index.html';
    });

    // Загрузка фильмов для управления
    async function loadMovies() {
        try {
            const token = getToken();
            const searchParams = new URLSearchParams();

            const response = await fetch(`${API_BASE}/movies?${searchParams}`, {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });

            if (response.status === 401) {
                removeToken();
                redirectToLogin();
                return;
            }

            const movies = await response.json();
            displayMoviesManagement(movies);
        } catch (error) {
            showMessage('Ошибка загрузки фильмов', 'error');
        }
    }

    // Отображение фильмов с кнопками управления
    function displayMoviesManagement(movies) {
        const container = document.getElementById('moviesContainer');

        container.innerHTML = '';

        if (movies.length === 0) {
            container.innerHTML = '<p>Фильмы не найдены</p>';
            return;
        }

        container.innerHTML = movies.map(movie => `
            <div class="movie-card">
                <h3>${movie.title}</h3>
                <p><strong>Жанр:</strong> ${movie.genre}</p>
                <p><strong>Режиссер:</strong> ${movie.director}</p>
                <p><strong>Рейтинг:</strong> <span class="rating">${movie.rating}</span>/10</p>
                <p><strong>Дата выхода:</strong> ${new Date(movie.release_date).toLocaleDateString('ru-RU')}</p>
                <div class="movie-actions">
                    <button class="btn" onclick="editMovie(${movie.id})">Редактировать</button>
                    <button class="btn btn-danger" onclick="deleteMovie(${movie.id})">Удалить</button>
                </div>
            </div>
        `).join('');
    }

    // Добавление/редактирование фильма
    document.getElementById('movieForm').addEventListener('submit', async (e) => {
        e.preventDefault();

        const formData = {
            title: document.getElementById('title').value,
            genre: document.getElementById('genre').value,
            director: document.getElementById('director').value,
            rating: parseFloat(document.getElementById('rating').value),
            release_date: document.getElementById('release_date').value
        };

        try {
            const token = getToken();
            let response;

            if (currentEditId) {
                // Редактирование
                response = await fetch(`${API_BASE}/movies/${currentEditId}`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${token}`
                    },
                    body: JSON.stringify(formData)
                });
            } else {
                // Добавление
                response = await fetch(`${API_BASE}/movies`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${token}`
                    },
                    body: JSON.stringify(formData)
                });
            }

            // В блоке добавления/редактирования фильма замените обработчик ошибок:
            if (response.ok) {
                const result = await response.json();
                console.log('Success response:', result);
                showMessage(currentEditId ? 'Фильм обновлен!' : 'Фильм добавлен!', 'success');
                resetForm();
                loadMovies();
            } else {
                // Получаем детальную информацию об ошибке
                let errorMessage = 'Ошибка сохранения фильма';
                try {
                    const errorData = await response.json();
                    errorMessage = errorData.error || errorMessage;
                } catch (e) {
                    errorMessage = `${response.status} ${response.statusText}`;
                }
                console.error('Error response:', errorMessage);
                showMessage(errorMessage, 'error');
            }
        } catch (error) {
            showMessage('Ошибка сети', 'error');
        }
    });

    // Редактирование фильма
    window.editMovie = async (id) => {
        console.log('Editing movie with ID:', id);
        try {
            const token = getToken();
            const response = await fetch(`${API_BASE}/movies`, {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const movies = await response.json();
            const movie = movies.find(m => m.id === id);

            if (movie) {
                document.getElementById('formTitle').textContent = 'Редактировать фильм';

                // Заполняем форму данными фильма
                document.getElementById('title').value = movie.title;
                document.getElementById('genre').value = movie.genre;
                document.getElementById('director').value = movie.director;
                document.getElementById('rating').value = movie.rating;

                // Правильно форматируем дату для input[type=date]
                let releaseDate;
                if (movie.release_date) {
                    if (typeof movie.release_date === 'string') {
                        // Если дата пришла как строка
                        if (movie.release_date.includes('T')) {
                            releaseDate = movie.release_date.split('T')[0];
                        } else {
                            releaseDate = movie.release_date;
                        }
                    } else {
                        // Если это объект Date (маловероятно в JSON)
                        releaseDate = new Date(movie.release_date).toISOString().split('T')[0];
                    }
                }
                document.getElementById('releaseDate').value = releaseDate || '';

                currentEditId = movie.id;
                document.getElementById('cancelEdit').style.display = 'inline-block';

                console.log('Form filled with movie data for editing:', movie);
                console.log('currentEditId set to:', currentEditId);

                // Прокрутка к форме
                document.querySelector('.movie-form').scrollIntoView({ behavior: 'smooth' });
            } else {
                showMessage('Фильм не найден', 'error');
            }
        } catch (error) {
            console.error('Error loading movie for edit:', error);
            showMessage('Ошибка загрузки фильма: ' + error.message, 'error');
        }
    };
    // Удаление фильма
    window.deleteMovie = async (id) => {
        if (!confirm('Вы уверены, что хотите удалить этот фильм?')) {
            return;
        }

        try {
            const token = getToken();
            const response = await fetch(`${API_BASE}/movies/${id}`, {
                method: 'DELETE',
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });

            if (response.ok) {
                showMessage('Фильм удален!', 'success');
                loadMovies();
            } else {
                showMessage('Ошибка удаления фильма', 'error');
            }
        } catch (error) {
            showMessage('Ошибка сети', 'error');
        }
    };

    // Сброс формы
    function resetForm() {
        document.getElementById('movieForm').reset();
        document.getElementById('formTitle').textContent = 'Добавить новый фильм';
        document.getElementById('movieId').value = '';
        currentEditId = null;
        document.getElementById('cancelEdit').style.display = 'none';
    }


    async function checkMovieExists(id) {
        try {
            const token = getToken();
            const response = await fetch(`${API_BASE}/movies`, {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });

            if (!response.ok) return false;

            const movies = await response.json();
            return movies.some(movie => movie.id === id);
        } catch (error) {
            console.error('Error checking movie existence:', error);
            return false;
        }
    }

    // Отмена редактирования
    document.getElementById('cancelEdit').addEventListener('click', resetForm);

    // Загрузка фильмов при загрузке страницы
    document.addEventListener('DOMContentLoaded', loadMovies);
}