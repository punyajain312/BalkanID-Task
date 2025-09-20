import { useState } from "react";
import { signup, login } from "../api/auth";
import SharedInput from "./SharedInput";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import toast from "react-hot-toast";

export default function SignupForm() {
  const { login: setAuth } = useAuth();
  const navigate = useNavigate();

  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await signup(name, email, password);
      const res = await login(email, password);
      setAuth(res.data.token);
      toast.success("Signup successful! Redirecting to dashboard...");
      navigate("/dashboard");
    } catch (err: any) {
      toast.error("Signup failed");
      console.error(err);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <h2 className="text-xl font-bold">Signup</h2>
      <SharedInput type="text" placeholder="Name" value={name} onChange={setName} label="Full Name" />
      <SharedInput type="email" placeholder="Email" value={email} onChange={setEmail} label="Email" />
      <SharedInput type="password" placeholder="Password" value={password} onChange={setPassword} label="Password" />
      <button type="submit" className="w-full bg-green-500 text-white p-2 rounded">
        Sign Up
      </button>
    </form>
  );
}