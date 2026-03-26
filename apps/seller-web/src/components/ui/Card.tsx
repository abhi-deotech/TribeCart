
import React from 'react';

interface CardProps {
  title?: string;
  children: React.ReactNode;
  className?: string;
}

const Card: React.FC<CardProps> = ({ title, children, className = '' }) => {
  return (
    <div className={`bg-white rounded-lg shadow-md p-5 ${className}`}>
      {title && (
        <h3 className="text-xl font-bold text-gray-800 mb-3 text-center md:text-left">
          {title}
        </h3>
      )}
      <div>{children}</div>
    </div>
  );
};

export default Card;
