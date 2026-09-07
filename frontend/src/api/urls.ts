import client from "./client";
import type {
	ShortURL,
	URLStats,
	CreateURLRequest,
	UpdateURLRequest,
} from "../types";

export const getAllURLs = async (): Promise<ShortURL[]> => {
	const res = await client.get<ShortURL[]>("/shorten");
	return res.data;
};

export const getURLByCode = async (code: string): Promise<ShortURL> => {
	const res = await client.get<ShortURL>(`/shorten/${code}`);
	return res.data;
};

export const createURL = async (data: CreateURLRequest): Promise<ShortURL> => {
	const res = await client.post<ShortURL>("/shorten", data);
	return res.data;
};

export const updateURL = async (
	code: string,
	data: UpdateURLRequest,
): Promise<ShortURL> => {
	const res = await client.put<ShortURL>(`/shorten/${code}`, data);
	return res.data;
};

export const deleteURL = async (code: string): Promise<void> => {
	await client.delete(`/shorten/${code}`);
};

export const getURLStats = async (code: string): Promise<URLStats> => {
	const res = await client.get<URLStats>(`/shorten/${code}/stats`);
	return res.data;
};
