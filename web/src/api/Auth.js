import React, { useEffect, useState } from "react";
import { useHistory } from "react-router-dom";
import { Request } from "./Base";
import Roles from "./Roles";

const AuthContext = React.createContext({ user: null });

const GetCurrentUser = async () => {
  const res = await Request("GET", "/whoami");

  if (res.error !== null) {
    return null;
  }
  return res.data;
};

const DoLogin = async (email, password) =>
  await Request("POST", "/login", { email, password });

const DoLogout = async () => await Request("POST", "/logout");

const AuthProvider = ({ children }) => {
  const history = useHistory();

  const [user, setUser] = useState(null);
  const [lastCheck, setLastCheck] = useState(Date.now());

  useEffect(() => {
    (async () => {
      const currentUser = await GetCurrentUser();
      setUser(currentUser);
    })();
  }, [lastCheck]);

  const logout = () => {
    setUser(null);
    DoLogout();

    // Return to the login page.
    history.push("/login");
  };

  const refreshLoginStatus = () => {
    setLastCheck(Date.now());
  };

  const isAuthenticated = () => user !== null;

  const getRole = () => user?.["role_id"] ?? false;

  const isOneOfRoles = (roles) => roles.includes(getRole());

  const isRole = (role) => getRole() === role;

  const isAdmin = () => isRole(Roles.IDOf.ADMIN);
  const isSponsor = () => isRole(Roles.IDOf.SPONSOR);
  const isDriver = () => isRole(Roles.IDOf.DRIVER);

  const getName = () => {
    const firstName = user?.["first_name"] ?? "Unknown";
    const lastName = user?.["last_name"] ?? "Unknown";

    return [firstName, lastName];
  };

  const getInitials = () => {
    const [firstName, lastName] = getName();
    const initials = (
      firstName.substring(0, 1) + lastName.substring(0, 1)
    ).toUpperCase();

    return initials;
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        logout,

        getRole,
        getName,
        getInitials,

        isAuthenticated,
        isOneOfRoles,
        isAdmin,
        isSponsor,
        isDriver,

        refreshLoginStatus,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

const WithUser = AuthContext.Consumer;

export { AuthProvider, WithUser, DoLogin, DoLogout, GetCurrentUser };
