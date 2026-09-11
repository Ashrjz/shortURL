import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { URLTable } from "../components/URLTable";
import { CreateURLForm } from "../components/CreateURLForm";
import * as urlsApi from "../api/urls";
import type { ShortURL } from "../types";

export default function Dashboard() {
	const [urls, setUrls] = useState<ShortURL[]>([]);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState("");

	const { logout } = useAuth();
	const navigate = useNavigate();

	const fetchUrls = async () => {
		try {
			setLoading(true);
			const data = await urlsApi.getAllURLs();
			setUrls(data);
		} catch (err) {
			setError("Failed to load URLs");
		} finally {
			setLoading(false);
		}
	};

	useEffect(() => {
		fetchUrls();
	}, []);

	const handleDelete = async (code: string) => {
		if (!confirm(`Delete short URL "${code}"?`)) return;
		try {
			await urlsApi.deleteURL(code);
			setUrls((prev) => prev.filter((u) => u.shortCode !== code));
		} catch (err) {
			setError("Failed to delete URL");
		}
	};

	const handleLogout = () => {
		logout();
		navigate("/login");
	};

	return (
		<div className="min-h-screen bg-neutral-950">
			<header className="bg-neutral-900 border-b border-neutral-800">
				<div className="max-w-5xl mx-auto px-4 py-4 flex justify-between items-center">
					<h1 className="text-xl font-bold text-neutral-100">
						URL Shortener{" "}
						<span className="text-emerald-400">Dashboard</span>
					</h1>
					<button
						onClick={handleLogout}
						className="text-sm text-neutral-400 hover:text-red-400"
					>
						Logout
					</button>
				</div>
			</header>

			<main className="max-w-5xl mx-auto px-4 py-8">
				{error && (
					<div className="mb-4 text-sm text-red-400 bg-red-950/40 border border-red-900 p-3 rounded">
						{error}
					</div>
				)}

				<CreateURLForm
					onCreated={(newUrl) => setUrls((prev) => [newUrl, ...prev])}
				/>

				{loading ? (
					<div className="text-center py-12 text-neutral-500">
						Loading...
					</div>
				) : (
					<URLTable urls={urls} onDelete={handleDelete} />
				)}
			</main>
		</div>
	);
}
