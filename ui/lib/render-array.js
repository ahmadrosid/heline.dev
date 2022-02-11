import React from "react";

export default function renderArray(child) {
    return React.Children.toArray(child)
}