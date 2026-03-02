import { useEffect, useState, type CSSProperties } from "react";

const UNITS = {
  length: ["millimeter", "centimeter", "meter", "kilometer", "inch", "foot", "yard", "mile"],
  weight: ["milligram", "gram", "kilogram", "ounce", "pound"],
  temperature: ["celsius", "fahrenheit", "kelvin"],
} as const;

type UnitType = keyof typeof UNITS; // "length" | "weight" | "temperature"
const TYPES = ["length", "weight", "temperature"] as const;

type ConvertOk = { result: number };
type ConvertErr = { error: string };

export default function App() {
  const [type, setType] = useState<UnitType>("length");
  const [value, setValue] = useState<string>("");

  const units = UNITS[type]; // readonly string[]
  const [from, setFrom] = useState<string>(units[0]);
  const [to, setTo] = useState<string>(units[1]);

  const [result, setResult] = useState<number | null>(null);
  const [error, setError] = useState<string>("");

  useEffect(() => {
    setFrom(UNITS[type][0]);
    setTo(UNITS[type][1]);
    setResult(null);
    setError("");
    setValue("");
  }, [type]);

  async function onConvert() {
    setError("");
    setResult(null);

    const num = Number(value);
    if (!Number.isFinite(num)) {
      setError("value must be a number");
      return;
    }

    try {
      const res = await fetch("http://localhost:8080/convert", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ type, value: num, from, to }),
      });

      const data = (await res.json()) as Partial<ConvertOk & ConvertErr>;

      if (!res.ok) {
        setError(data.error ?? "request failed");
        return;
      }

      setResult(typeof data.result === "number" ? data.result : null);
    } catch {
      setError("cannot reach backend (is Go server running on :8080?)");
    }
  }

  function onReset() {
    setValue("");
    setResult(null);
    setError("");
    setFrom(UNITS[type][0]);
    setTo(UNITS[type][1]);
  }

  const containerStyle: CSSProperties = {
    maxWidth: 520,
    margin: "40px auto",
    padding: "0 16px",
    fontFamily: "sans-serif",
    transform: "translateX(30px)", // <-- 오른쪽으로 살짝 이동
  };

  return (
    <div style={containerStyle}>
      <h1>Unit Converter</h1>

      <div style={{ display: "flex", gap: 10, marginBottom: 20 }}>
        {TYPES.map((t) => (
          <button
            key={t}
            onClick={() => setType(t)}
            style={{
              padding: "8px 12px",
              fontWeight: type === t ? "700" : "400",
            }}
          >
            {t}
          </button>
        ))}
      </div>

      <div style={{ display: "grid", gap: 12 }}>
        <label>
          Enter the value to convert
          <input
            value={value}
            onChange={(e) => setValue(e.target.value)}
            placeholder="e.g. 123"
            style={{ width: "100%", padding: 8, marginTop: 6 }}
          />
        </label>

        <label>
          Unit to convert from
          <select
            value={from}
            onChange={(e) => setFrom(e.target.value)}
            style={{ width: "100%", padding: 8, marginTop: 6 }}
          >
            {units.map((u) => (
              <option key={u} value={u}>
                {u}
              </option>
            ))}
          </select>
        </label>

        <label>
          Unit to convert to
          <select
            value={to}
            onChange={(e) => setTo(e.target.value)}
            style={{ width: "100%", padding: 8, marginTop: 6 }}
          >
            {units.map((u) => (
              <option key={u} value={u}>
                {u}
              </option>
            ))}
          </select>
        </label>

        <div style={{ display: "flex", gap: 10 }}>
          <button onClick={onConvert} style={{ padding: "10px 14px" }}>
            Convert
          </button>
          <button onClick={onReset} style={{ padding: "10px 14px" }}>
            Reset
          </button>
        </div>

        {error && (
          <div style={{ padding: 10, border: "1px solid #f99" }}>
            <b>Error:</b> {error}
          </div>
        )}

        {result !== null && (
          <div style={{ padding: 10, border: "1px solid #999" }}>
            <div>Result</div>
            <div style={{ fontSize: 28, fontWeight: 700 }}>
              {value} {from} = {result} {to}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}