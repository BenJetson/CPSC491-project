import React from "react";
import { Route, Switch, useRouteMatch } from "react-router-dom";

import NotImplemented from "./components/NotImplemented";
import NotFound from "./components/NotFound";

import { WithUser } from "./api/Auth";
import Roles from "./api/Roles";
import AccessDenied from "./components/AccessDenied";
import VendorSearch from "./components/VendorSearch";
import SponsorCatalog from "./components/SponsorCatalog";
import OrgProfileEditor from "./components/OrganizationEditor";

const AppSubrouterSponsor = () => {
  const match = useRouteMatch();

  return (
    <WithUser>
      {({ isOneOfRoles }) =>
        (isOneOfRoles([Roles.IDOf.SPONSOR]) && (
          <Switch>
            <Route path={`${match.path}/orgmgt`}>
              <OrgProfileEditor />
            </Route>
            <Route path={`${match.path}/applications`}>
              <NotImplemented feature={"Sponsor - Manage Applications"} />
            </Route>
            <Route path={`${match.path}/drivers`}>
              <NotImplemented feature={"Sponsor - Manage Drivers"} />
            </Route>
            <Route exact path={`${match.path}/catalog`}>
              <SponsorCatalog />
            </Route>
            <Route exact path={`${match.path}/catalog/vendor`}>
              <VendorSearch />
            </Route>
            <Route path={"*"}>
              {/* If no route matches, show a not found page. */}
              <NotFound />
            </Route>
          </Switch>
        )) || <AccessDenied />
      }
    </WithUser>
  );
};

export default AppSubrouterSponsor;
