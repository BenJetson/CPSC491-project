import React from "react";

import {
  Home as HomeIcon,
  PlayArrow as PlayArrowIcon,
  People as PeopleIcon,
  Business as BusinessIcon,
  Assignment as AssignmentIcon,
  LocalAtm as LocalAtmIcon,
  ShoppingBasket as ShoppingBasketIcon,
  Receipt as ReceiptIcon,
  Storefront as StorefrontIcon,
  AssignmentTurnedIn as AssignmentTurnedInIcon,
  Face as FaceIcon,
  AccountBox as AccountBoxIcon,
  Help as HelpIcon,
  Info as InfoIcon,
} from "@material-ui/icons";

import Roles from "../api/Roles";

let Navigation = [
  {
    title: null,
    roles: null,
    items: [
      {
        link: "/",
        name: "Home",
        icon: <HomeIcon />,
      },
      {
        link: "/help",
        name: "Help Center",
        icon: <HelpIcon />,
      },
      {
        link: "/about",
        name: "About This App",
        icon: <InfoIcon />,
      },
    ],
  },
  {
    title: "My Tools",
    roles: [Roles.ADMIN, Roles.SPONSOR, Roles.DRIVER],
    items: [
      {
        link: "/profile",
        name: "My Profile",
        icon: <AccountBoxIcon />,
      },
    ],
  },
  {
    title: "Administrative Tools",
    roles: [Roles.ADMIN],
    items: [
      {
        link: "/admin/users",
        name: "Manage Users",
        icon: <PeopleIcon />,
      },
      {
        link: "/admin/organizations",
        name: "Manage Organizations",
        icon: <BusinessIcon />,
      },
    ],
  },
  {
    title: "Sponsor Tools",
    roles: [Roles.ADMIN, Roles.SPONSOR],
    items: [
      {
        link: "/sponsor/applications",
        name: "Manage Applications",
        icon: <AssignmentTurnedInIcon />,
      },
      {
        link: "/sponsor/drivers",
        name: "Manage Drivers",
        icon: <FaceIcon />,
      },
      {
        link: "/sponsor/catalog",
        name: "Manage Catalog",
        icon: <StorefrontIcon />,
      },
    ],
  },
  {
    title: "Driver Tools",
    roles: [Roles.ADMIN, Roles.DRIVER],
    items: [
      {
        link: "/driver/apply",
        name: "Applications",
        icon: <AssignmentIcon />,
      },
      {
        link: "/driver/balance",
        name: "View Balance",
        icon: <LocalAtmIcon />,
      },
      {
        link: "/driver/shop",
        name: "Incentive Shop",
        icon: <ShoppingBasketIcon />,
      },
      {
        link: "/driver/orders",
        name: "Receipts",
        icon: <ReceiptIcon />,
      },
    ],
  },
];

export default Navigation;
