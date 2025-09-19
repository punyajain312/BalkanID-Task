import { useState } from "react";
import { login } from "../api/auth";
import { useAuth } from "../context/AuthContext";
import SharedInput from "./SharedInput";
import { useNavigate } from "react-router-dom";

export default function LoginForm() {
  const { login: setAuth } = useAuth();
  const navigate = useNavigate();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const res = await login(email, password);
      setAuth(res.data.token);
      navigate("/upload");
    } catch (err: any) {
      console.error(err);
      alert("Login failed. Check console for details.");
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <h2 className="text-xl font-bold">Login</h2>
      <SharedInput
        type="email"
        placeholder="Email"
        value={email}
        onChange={setEmail}
        label="Email Address"
      />
      <SharedInput
        type="password"
        placeholder="Password"
        value={password}
        onChange={setPassword}
        label="Password"
      />
      <button
        type="submit"
        className="w-full bg-blue-500 text-white p-2 rounded"
      >
        Login
      </button>
    </form>
  );
}