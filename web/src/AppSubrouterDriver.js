import React from "react";
import { Route, Switch, useRouteMatch } from "react-router-dom";
import NotFound from "./components/NotFound";
import NotImplemented from "./components/NotImplemented";

const AppSubrouterDriver = () => {
  const match = useRouteMatch();

  return (
    <Switch>
      <Route path={`${match.path}/applications`}>
        <NotImplemented feature={"Driver - Applications"} />
      </Route>
      <Route path={`${match.path}/balance`}>
        <NotImplemented feature={"Driver - View Balance"} />
      </Route>
      <Route path={`${match.path}/shop`}>
        <NotImplemented feature={"Driver - Incentive Shop"} />
      </Route>
      <Route path={`${match.path}/receipts`}>
        <NotImplemented feature={"Driver - Receipts"} />
      </Route>
      <Route path={"*"}>
        {/* If no route matches, show a not found page. */}
        <NotFound />
      </Route>
    </Switch>
  );
};

export default AppSubrouterDriver;
