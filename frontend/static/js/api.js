

const API = {

    async request(url, options = {}, requiresAuth = true) {
        const headers = {
            'Content-Type': 'application/json',
            ...options.headers
        };


        if (requiresAuth) {
            const accessToken = localStorage.getItem('access_token');
            if (!accessToken) {
                console.error('Authentication required but no access token found');
                this.redirectToLogin();
                return Promise.reject(new Error('Authentication required'));
            }
            headers['Authorization'] = `Bearer ${accessToken}`;
        }

        const fetchOptions = {
            ...options,
            headers
        };

        try {

            const response = await fetch(url, fetchOptions);

            if (response.status === 401 && requiresAuth) {
                console.log('Token expired, attempting to refresh...');
                const refreshed = await this.refreshToken();

                if (refreshed) {
                    headers['Authorization'] = `Bearer ${localStorage.getItem('access_token')}`;
                    return this.request(url, options, requiresAuth);
                } else {

                    this.redirectToLogin();
                    return Promise.reject(new Error('Authentication failed'));
                }
            }


            const data = await response.json();

            if (!response.ok) {
                console.error(`API Error: ${data.message || 'Unknown error'}`);
                return Promise.reject(data);
            }

            return data;
        } catch (error) {
            console.error('API Request failed:', error);
            return Promise.reject(error);
        }
    },


    async refreshToken() {
        const refreshToken = localStorage.getItem('refresh_token');

        if (!refreshToken) {
            console.error('No refresh token available');
            return false;
        }

        try {
            const response = await fetch('/api/auth/refresh', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ refresh_token: refreshToken })
            });

            if (!response.ok) {
                throw new Error('Failed to refresh token');
            }

            const data = await response.json();

            if (data.code === 200) {
                localStorage.setItem('access_token', data.access_token);
                if (data.refresh_token) {
                    localStorage.setItem('refresh_token', data.refresh_token);
                }
                return true;
            } else {
                throw new Error(data.message || 'Failed to refresh token');
            }
        } catch (error) {
            console.error('Token refresh failed:', error);
            return false;
        }
    },


    redirectToLogin() {
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        window.location.href = '/';
    },


    async get(url, requiresAuth = true) {
        return this.request(url, { method: 'GET' }, requiresAuth);
    },


    async post(url, data, requiresAuth = true) {
        return this.request(
            url,
            {
                method: 'POST',
                body: JSON.stringify(data)
            },
            requiresAuth
        );
    },


    async put(url, data, requiresAuth = true) {
        return this.request(
            url,
            {
                method: 'PUT',
                body: JSON.stringify(data)
            },
            requiresAuth
        );
    },


    async delete(url, data = null, requiresAuth = true) {
        const options = {
            method: 'DELETE'
        };

        if (data) {
            options.body = JSON.stringify(data);
        }

        return this.request(url, options, requiresAuth);
    },


    async login(email, password) {
        const data = await this.post('/api/auth/login', { email, password }, false);

        if (data.code === 200) {
            localStorage.setItem('access_token', data.access_token);
            localStorage.setItem('refresh_token', data.refresh_token);

            if (data.user.role === 'admin') {
                window.location.href = '/admin_main_menu';
            } else {
                window.location.href = '/personal_account';
            }
        }

        return data;
    },

    async logout() {
        try {
            const data = await this.post('/api/auth/logout', {});

            if (data.code === 200) {
                localStorage.removeItem('access_token');
                localStorage.removeItem('refresh_token');
                window.location.href = '/';
            }

            return data;
        } catch (error) {
            localStorage.removeItem('access_token');
            localStorage.removeItem('refresh_token');
            window.location.href = '/';
            throw error;
        }
    },

    async resetPassword(email) {
        return this.post('/api/auth/reset_password', { email }, false);
    }
};

export default API;