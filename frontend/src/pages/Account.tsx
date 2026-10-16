import { useAuth } from "../auth/useAuth";

export default function Account() {
  const user = useAuth((s) => s.user);

  return (
    <div className="narrow">
      <div className="section-head">
        <h1>Mi cuenta</h1>
      </div>
      <section className="card">
        <div className="field">
          <span className="label">Nombre</span>
          <span className="value">{user?.full_name}</span>
        </div>
        <div className="field">
          <span className="label">Email</span>
          <span className="value">{user?.email}</span>
        </div>
        <div className="field">
          <span className="label">Tipo de cuenta</span>
          <span className="value">{user?.role === "admin" ? "Administrador" : "Cliente"}</span>
        </div>
      </section>
    </div>
  );
}
