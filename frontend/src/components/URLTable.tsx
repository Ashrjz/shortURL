import { useState } from "react";
import { Link } from "react-router-dom";
import type { ShortURL } from "../types";

interface URLTableProps {
	urls: ShortURL[];
	onDelete: (code: string) => void;
}

export function URLTable({ urls, onDelete }: URLTableProps) {
	const [copiedCode, setCopiedCode] = useState<string | null>(null);

	const handleCopy = (code: string) => {
		const shortUrl = `http://localhost:8080/${code}`;
		navigator.clipboard.writeText(shortUrl);
		setCopiedCode(code);
		setTimeout(() => setCopiedCode(null), 1500);
	};

	if (urls.length === 0) {
		return (
			<div className="text-center py-12 text-neutral-500 bg-neutral-900 border border-neutral-800 rounded-lg">
				No URLs yet. Create your first short URL above.
			</div>
		);
	}

	return (
		<div className="overflow-x-auto bg-neutral-900 border border-neutral-800 rounded-lg shadow">
			<table className="min-w-full divide-y divide-neutral-800">
				<thead className="bg-neutral-800/50">
					<tr>
						<th className="px-4 py-3 text-center text-xs font-medium text-neutral-400 uppercase">
							Short Code
						</th>
						<th className="px-4 py-3 text-center text-xs font-medium text-neutral-400 uppercase">
							Original URL
						</th>
						<th className="px-4 py-3 text-center text-xs font-medium text-neutral-400 uppercase">
							Created
						</th>
						<th className="px-4 py-3 text-center text-xs font-medium text-neutral-400 uppercase">
							Actions
						</th>
					</tr>
				</thead>
				<tbody className="divide-y divide-neutral-800">
					{urls.map((u) => (
						<tr
							key={u.shortCode}
							className="hover:bg-neutral-800/30"
						>
							<td className="px-4 py-3 text-sm font-mono text-emerald-400">
								<Link
									to={`/urls/${u.shortCode}`}
									className="hover:text-emerald-300"
								>
									{u.shortCode}
								</Link>
							</td>
							<td className="px-4 py-3 text-sm text-neutral-300 max-w-xs truncate">
								{u.url}
							</td>
							<td className="px-4 py-3 text-sm text-neutral-500">
								{new Date(u.createdAt).toLocaleDateString()}
							</td>
							<td className="px-4 py-3 text-sm space-x-3">
								<button
									onClick={() => handleCopy(u.shortCode)}
									className="text-neutral-400 hover:text-emerald-400"
								>
									{copiedCode === u.shortCode
										? "Copied!"
										: "Copy"}
								</button>
								<Link
									to={`/urls/${u.shortCode}`}
									className="text-neutral-400 hover:text-emerald-400"
								>
									View
								</Link>
								<button
									onClick={() => onDelete(u.shortCode)}
									className="text-neutral-400 hover:text-red-400"
								>
									Delete
								</button>
							</td>
						</tr>
					))}
				</tbody>
			</table>
		</div>
	);
}
