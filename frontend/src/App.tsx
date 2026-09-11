import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { ProtectedRoute } from "./components/ProtectedRoute";
import Login from "./pages/Login";
import Register from "./pages/Register";
import Dashboard from "./pages/Dashboard";
import URLDetail from "./pages/URLDetail";

function App() {
	return (
		<BrowserRouter>
			<Routes>
				<Route path="/login" element={<Login />} />
				<Route path="/register" element={<Register />} />
				<Route
					path="/dashboard"
					element={
						<ProtectedRoute>
							<Dashboard />
						</ProtectedRoute>
					}
				/>
				<Route
					path="/urls/:code"
					element={
						<ProtectedRoute>
							<URLDetail />
						</ProtectedRoute>
					}
				/>
				<Route
					path="/"
					element={<Navigate to="/dashboard" replace />}
				/>
			</Routes>
		</BrowserRouter>
	);
}

export default App;
