import { useState } from "react";
import { useNavigate } from "react-router-dom";
import OrdersSection from "../components/OrdersSection";
import ProductsSection from "../components/ProductsSection";
import { useAuth } from "../auth/useAuth";

export default function Dashboard() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  // Al crear un pedido, cambiamos este valor para que "Mis pedidos" se recargue.
  const [ordersVersion, setOrdersVersion] = useState(0);

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
        <p>Elegí productos, armá tu pedido y pagalo desde acá.</p>
      </section>

      <section>
        <h3>Catálogo</h3>
        <ProductsSection onOrderCreated={() => setOrdersVersion((v) => v + 1)} />
      </section>

      <section>
        <h3>Mis pedidos</h3>
        <OrdersSection refreshKey={ordersVersion} />
      </section>

      <section>
        <h3>Datos de tu cuenta</h3>
        <div className="field">
          <span className="label">Email</span>
          <span className="value">{user?.email}</span>
        </div>
        <div className="field">
          <span className="label">Nombre</span>
          <span className="value">{user?.full_name}</span>
        </div>
        <div className="field">
          <span className="label">Rol</span>
          <span className="value">{user?.role}</span>
        </div>
      </section>
    </div>
  );
}
