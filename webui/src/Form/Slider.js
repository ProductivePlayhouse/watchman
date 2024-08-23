import React from 'react';
import styled from "styled-components";
import MSlider from "@mui/material/Slider";
import Typography from "@mui/material/Typography";

const render = ({ label, ...props }) => (
  <div
    css={`
      width: 100%;
    `}
  >
    <Typography
      css={`
        && {
        color: rgba(0, 0, 0, 0.54);
      `}
    >
      {label}
    </Typography>
    <MSlider {...props} />
  </div>
);

export default render;
