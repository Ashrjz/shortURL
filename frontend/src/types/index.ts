export interface ShortURL {
	id: number;
	url: string;
	shortCode: string;
	createdAt: string;
	updatedAt: string;
}

export interface URLStats extends ShortURL {
	accessCount: number;
}

export interface AuthResponse {
	token: string;
}

export interface RegisterRequest {
	username: string;
	password: string;
}

export interface LoginRequest {
	username: string;
	password: string;
}

export interface CreateURLRequest {
	url: string;
}

export interface UpdateURLRequest {
	url: string;
}

export interface ApiError {
	error: string;
}
