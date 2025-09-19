interface Props {
  type: string;
  placeholder: string;
  value: string;
  onChange: (val: string) => void;
}

export default function SharedInput({ type, placeholder, value, onChange }: Props) {
  return (
    <input
      type={type}
      placeholder={placeholder}
      value={value}
      onChange={(e) => onChange(e.target.value)}
      className="w-full p-2 border rounded"
      required
    />
  );
}