import { createContext, useContext, useState, useEffect } from "react";
import type { ReactNode } from "react";
import * as authApi from "../api/auth";
import type { LoginRequest, RegisterRequest } from "../types";

interface AuthContextType {
	token: string | null;
	isAuthenticated: boolean;
	login: (data: LoginRequest) => Promise<void>;
	register: (data: RegisterRequest) => Promise<void>;
	logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
	const [token, setToken] = useState<string | null>(
		localStorage.getItem("token"),
	);

	useEffect(() => {
		if (token) {
			localStorage.setItem("token", token);
		} else {
			localStorage.removeItem("token");
		}
	}, [token]);

	const login = async (data: LoginRequest) => {
		const res = await authApi.login(data);
		setToken(res.token);
	};

	const register = async (data: RegisterRequest) => {
		const res = await authApi.register(data);
		setToken(res.token);
	};

	const logout = () => {
		setToken(null);
	};

	return (
		<AuthContext.Provider
			value={{ token, isAuthenticated: !!token, login, register, logout }}
		>
			{children}
		</AuthContext.Provider>
	);
}

export function useAuth() {
	const context = useContext(AuthContext);
	if (!context) {
		throw new Error("useAuth must be used within AuthProvider");
	}
	return context;
}
