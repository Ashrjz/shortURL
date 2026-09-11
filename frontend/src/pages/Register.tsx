import { useState, type SyntheticEvent } from "react";
import { useNavigate, Link } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

export default function Register() {
	const [username, setUsername] = useState("");
	const [password, setPassword] = useState("");
	const [error, setError] = useState("");
	const [loading, setLoading] = useState(false);

	const { register } = useAuth();
	const navigate = useNavigate();

	const handleSubmit = async (e: SyntheticEvent) => {
		e.preventDefault();
		setError("");
		setLoading(true);
		try {
			await register({ username, password });
			navigate("/dashboard");
		} catch (err) {
			setError("Registration failed. Username may already exist.");
		} finally {
			setLoading(false);
		}
	};

	return (
		<div className="min-h-screen flex items-center justify-center bg-neutral-950">
			<form
				onSubmit={handleSubmit}
				className="bg-neutral-900 border border-neutral-800 p-8 rounded-lg shadow-md w-full max-w-sm"
			>
				<h1 className="text-2xl font-bold mb-6 text-neutral-100">
					Register
				</h1>

				{error && (
					<div className="mb-4 text-sm text-red-400 bg-red-950/40 border border-red-900 p-2 rounded">
						{error}
					</div>
				)}

				<div className="mb-4">
					<label className="block text-sm font-medium text-neutral-300 mb-1">
						Username
					</label>
					<input
						type="text"
						value={username}
						onChange={(e) => setUsername(e.target.value)}
						required
						className="w-full bg-neutral-800 border border-neutral-700 text-neutral-100 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-emerald-500"
					/>
				</div>

				<div className="mb-6">
					<label className="block text-sm font-medium text-neutral-300 mb-1">
						Password
					</label>
					<input
						type="password"
						value={password}
						onChange={(e) => setPassword(e.target.value)}
						required
						minLength={6}
						className="w-full bg-neutral-800 border border-neutral-700 text-neutral-100 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-emerald-500"
					/>
				</div>

				<button
					type="submit"
					disabled={loading}
					className="w-full bg-emerald-600 text-white py-2 rounded hover:bg-emerald-500 disabled:opacity-50 transition-colors"
				>
					{loading ? "Creating account..." : "Register"}
				</button>

				<p className="mt-4 text-sm text-neutral-400 text-center">
					Already have an account?{" "}
					<Link
						to="/login"
						className="text-emerald-400 hover:text-emerald-300"
					>
						Login
					</Link>
				</p>
			</form>
		</div>
	);
}
