import React from "react";
import styled from "styled-components";
import MTextField from "@mui/material/TextField";

// Styled Component
const StyledTextField = styled(MTextField)`
  && {
    min-width: 100%;
  }
`;

const Render = ({ id, ...props }) =>
{
  return <StyledTextField name={id} {...props} />;
};

export default Render;
