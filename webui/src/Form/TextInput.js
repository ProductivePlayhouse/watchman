import React from "react";
import styled from "styled-components";
import MTextField from "@mui/material/TextField";

const render = ({ id, ...props }) =>
{
  return (
    <MTextField
      css={`
        && {
          min-width: 100%;
        }
      `}
      name={id}
      {...props}
    />
  );
};

export default render;
