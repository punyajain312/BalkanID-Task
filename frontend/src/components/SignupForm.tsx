import { useState } from "react";
import { signup } from "../api/auth";
import SharedInput from "./SharedInput";

export default function SignupForm() {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await signup(name, email, password);
      alert("Signup successful! Please log in.");
    } catch (err: any) {
      alert("Signup failed: " + (err.response?.data?.error || err.message));
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <h2 className="text-xl font-bold">Signup</h2>
      <SharedInput type="text" placeholder="Name" value={name} onChange={setName} />
      <SharedInput type="email" placeholder="Email" value={email} onChange={setEmail} />
      <SharedInput type="password" placeholder="Password" value={password} onChange={setPassword} />
      <button type="submit" className="w-full bg-green-500 text-white p-2 rounded">
        Sign Up
      </button>
    </form>
  );
}