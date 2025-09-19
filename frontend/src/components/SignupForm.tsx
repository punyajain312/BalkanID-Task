import { useState } from "react";
import { signup } from "../api/auth";
import SharedInput from "./SharedInput";
import { useNavigate } from "react-router-dom";
import toast from "react-hot-toast";

export default function SignupForm() {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await signup(name, email, password);
      toast.success("Signup successful! Please log in 🔑");
      navigate("/login");
    } catch (err: any) {
      toast.error("Signup failed ❌");
      console.error(err);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <h2 className="text-xl font-bold">Signup</h2>
      <SharedInput
        type="text"
        placeholder="Name"
        value={name}
        onChange={setName}
        label="Full Name"
      />
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
        className="w-full bg-green-500 text-white p-2 rounded"
      >
        Sign Up
      </button>
    </form>
  );
}