import { useState, type SyntheticEvent } from "react";
import * as urlsApi from "../api/urls";
import type { ShortURL } from "../types";

interface CreateURLFormProps {
	onCreated: (url: ShortURL) => void;
}

export function CreateURLForm({ onCreated }: CreateURLFormProps) {
	const [url, setUrl] = useState("");
	const [error, setError] = useState("");
	const [loading, setLoading] = useState(false);

	const handleSubmit = async (e: SyntheticEvent) => {
		e.preventDefault();
		setError("");
		setLoading(true);
		try {
			const created = await urlsApi.createURL({ url });
			onCreated(created);
			setUrl("");
		} catch (err) {
			setError(
				"Failed to create short URL. Make sure it starts with http:// or https://",
			);
		} finally {
			setLoading(false);
		}
	};

	return (
		<form
			onSubmit={handleSubmit}
			className="bg-neutral-900 border border-neutral-800 p-4 rounded-lg shadow mb-6 flex gap-3 items-start"
		>
			<div className="flex-1">
				<input
					type="text"
					value={url}
					onChange={(e) => setUrl(e.target.value)}
					placeholder="https://example.com/long/url"
					required
					className="w-full bg-neutral-800 border border-neutral-700 text-neutral-100 placeholder-neutral-500 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-emerald-500"
				/>
				{error && <p className="text-sm text-red-400 mt-1">{error}</p>}
			</div>
			<button
				type="submit"
				disabled={loading}
				className="bg-emerald-600 text-white px-4 py-2 rounded hover:bg-emerald-500 disabled:opacity-50 whitespace-nowrap transition-colors"
			>
				{loading ? "Creating..." : "Shorten URL"}
			</button>
		</form>
	);
}
