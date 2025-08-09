import React from 'react';

interface SARSIconProps {
  width?: number;
  className?: string;
}

const SARSIcon: React.FC<SARSIconProps> = ({ 
  width = 104,
  className 
}) => {
  return (
<svg width={width} height="auto" viewBox="0 0 296 296" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M174 188.5V231L238.5 178.5V136.5L174 188.5Z" fill="#004CFF"/>
<path d="M238.5 95V52L148.5 125.5L57 52V211.5L94.5 243.5V125.5L148.5 168.5L238.5 95Z" fill="white"/>
</svg>
  );
};

export default SARSIcon;