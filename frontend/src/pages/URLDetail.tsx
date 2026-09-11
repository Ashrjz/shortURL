import { useEffect, useState } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import * as urlsApi from "../api/urls";
import type { URLStats } from "../types";

export default function URLDetail() {
	const { code } = useParams<{ code: string }>();
	const navigate = useNavigate();

	const [stats, setStats] = useState<URLStats | null>(null);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState("");

	const [editUrl, setEditUrl] = useState("");
	const [editing, setEditing] = useState(false);
	const [saving, setSaving] = useState(false);

	const fetchStats = async () => {
		if (!code) return;
		try {
			setLoading(true);
			const data = await urlsApi.getURLStats(code);
			setStats(data);
			setEditUrl(data.url);
		} catch (err) {
			setError("URL not found");
		} finally {
			setLoading(false);
		}
	};

	useEffect(() => {
		fetchStats();
	}, [code]);

	const handleUpdate = async () => {
		if (!code) return;
		setSaving(true);
		try {
			await urlsApi.updateURL(code, { url: editUrl });
			setEditing(false);
			fetchStats();
		} catch (err) {
			setError("Failed to update URL");
		} finally {
			setSaving(false);
		}
	};

	const handleDelete = async () => {
		if (!code || !confirm(`Delete short URL "${code}"?`)) return;
		try {
			await urlsApi.deleteURL(code);
			navigate("/dashboard");
		} catch (err) {
			setError("Failed to delete URL");
		}
	};

	const shortUrl = `http://localhost:8080/${code}`;

	if (loading) {
		return (
			<div className="min-h-screen bg-neutral-950 text-center py-12 text-neutral-500">
				Loading...
			</div>
		);
	}

	if (error || !stats) {
		return (
			<div className="min-h-screen flex flex-col items-center justify-center gap-4 bg-neutral-950">
				<p className="text-red-400">{error || "Not found"}</p>
				<Link
					to="/dashboard"
					className="text-emerald-400 hover:text-emerald-300"
				>
					Back to Dashboard
				</Link>
			</div>
		);
	}

	return (
		<div className="min-h-screen bg-neutral-950">
			<header className="bg-neutral-900 border-b border-neutral-800">
				<div className="max-w-3xl mx-auto px-4 py-4">
					<Link
						to="/dashboard"
						className="text-sm text-emerald-400 hover:text-emerald-300"
					>
						&larr; Back to Dashboard
					</Link>
				</div>
			</header>

			<main className="max-w-3xl mx-auto px-4 py-8">
				{error && (
					<div className="mb-4 text-sm text-red-400 bg-red-950/40 border border-red-900 p-3 rounded">
						{error}
					</div>
				)}

				<div className="bg-neutral-900 border border-neutral-800 rounded-lg shadow p-6 mb-6">
					<div className="flex justify-between items-start mb-4">
						<div>
							<p className="text-sm text-neutral-500">
								Short URL
							</p>
							<a
								href={shortUrl}
								target="_blank"
								rel="noopener noreferrer"
								className="text-lg font-mono text-emerald-400 hover:text-emerald-300"
							>
								{shortUrl}
							</a>
						</div>
						<button
							onClick={handleDelete}
							className="text-sm text-red-400 hover:text-red-300"
						>
							Delete
						</button>
					</div>

					<div className="mb-4">
						<p className="text-sm text-neutral-500 mb-1">
							Original URL
						</p>
						{editing ? (
							<div className="flex gap-2">
								<input
									type="text"
									value={editUrl}
									onChange={(e) => setEditUrl(e.target.value)}
									className="flex-1 bg-neutral-800 border border-neutral-700 text-neutral-100 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-emerald-500"
								/>
								<button
									onClick={handleUpdate}
									disabled={saving}
									className="bg-emerald-600 text-white px-4 py-2 rounded hover:bg-emerald-500 disabled:opacity-50 transition-colors"
								>
									{saving ? "Saving..." : "Save"}
								</button>
								<button
									onClick={() => {
										setEditing(false);
										setEditUrl(stats.url);
									}}
									className="text-neutral-400 px-3 py-2 hover:text-neutral-200"
								>
									Cancel
								</button>
							</div>
						) : (
							<div className="flex justify-between items-center">
								<p className="text-neutral-300 break-all">
									{stats.url}
								</p>
								<button
									onClick={() => setEditing(true)}
									className="text-sm text-emerald-400 hover:text-emerald-300 whitespace-nowrap ml-3"
								>
									Edit
								</button>
							</div>
						)}
					</div>

					<div className="grid grid-cols-3 gap-4 pt-4 border-t border-neutral-800 text-sm">
						<div>
							<p className="text-neutral-500">Access Count</p>
							<p className="text-xl font-bold text-emerald-400">
								{stats.accessCount}
							</p>
						</div>
						<div>
							<p className="text-neutral-500">Created</p>
							<p className="text-neutral-300">
								{new Date(stats.createdAt).toLocaleString()}
							</p>
						</div>
						<div>
							<p className="text-neutral-500">Last Updated</p>
							<p className="text-neutral-300">
								{new Date(stats.updatedAt).toLocaleString()}
							</p>
						</div>
					</div>
				</div>
			</main>
		</div>
	);
}
