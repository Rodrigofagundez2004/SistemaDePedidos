import { useNavigate } from "react-router-dom";
import { useAuth } from "../auth/useAuth";

export default function Dashboard() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate("/login");
  };

  return (
    <div className="dashboard-container">
      <header>
        <h1>Mi cuenta</h1>
        <button onClick={handleLogout}>Cerrar sesión</button>
      </header>

      <section>
        <h2>Hola, {user?.full_name}</h2>
        <p>Te damos la bienvenida a tu panel.</p>
      </section>

      <section>
        <h3>Datos de tu cuenta</h3>
        <div className="field">
          <span className="label"> Email: </span>
          <span className="value">{user?.email}</span>
        </div>
        <div className="field">
          <span className="label"> Nombre: </span>
          <span className="value">{user?.full_name}</span>
        </div>
        <div className="field">
          <span className="label"> Rol: </span>
          <span className="value">{user?.role}</span>
        </div>
      </section>

      <section>
        <h3>Próximamente</h3>
        <p>Acá vas a poder ver tus productos y gestionar tus órdenes.</p>
      </section>
    </div>
  );
}