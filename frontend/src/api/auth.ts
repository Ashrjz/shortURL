import client from "./client";
import type { AuthResponse, RegisterRequest, LoginRequest } from "../types";

export const register = async (
	data: RegisterRequest,
): Promise<AuthResponse> => {
	const res = await client.post<AuthResponse>("/register", data);
	return res.data;
};

export const login = async (data: LoginRequest): Promise<AuthResponse> => {
	const res = await client.post<AuthResponse>("/login", data);
	return res.data;
};
