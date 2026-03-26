
import React, { useState, useRef, useEffect } from 'react';
import { createPortal } from 'react-dom';
import { MoreVertical } from 'lucide-react';
import Button from './Button';

interface Action {
  label: string;
  icon?: React.ElementType;
  onClick: (event: React.MouseEvent) => void;
}

interface ActionDropdownProps {
  actions: Action[];
}

const ActionDropdown: React.FC<ActionDropdownProps> = ({ actions }) => {
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);
  const buttonRef = useRef<HTMLButtonElement>(null);
  const [dropdownStyle, setDropdownStyle] = useState({});

  const toggleDropdown = (e: React.MouseEvent) => {
    e.stopPropagation();
    setIsOpen((prev) => !prev);
  };

  const handleClickOutside = (event: MouseEvent) => {
    if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
      setIsOpen(false);
    }
  };

  useEffect(() => {
    document.addEventListener('click', handleClickOutside);
    return () => {
      document.removeEventListener('click', handleClickOutside);
    };
  }, []);

  useEffect(() => {
    if (isOpen && buttonRef.current) {
      const rect = buttonRef.current.getBoundingClientRect();
      const viewportWidth = window.innerWidth || document.documentElement.clientWidth;
      const viewportHeight = window.innerHeight || document.documentElement.clientHeight;

      let top = rect.bottom + window.scrollY;
      let left = rect.right + window.scrollX - 192; // 192px is w-48

      if (left + 192 > viewportWidth) {
        left = viewportWidth - 192 - 10;
      }
      if (left < 0) {
        left = 10;
      }

      if (top + 200 > viewportHeight + window.scrollY) {
        top = rect.top + window.scrollY - 200;
        if (top < window.scrollY) {
          top = window.scrollY + 10;
        }
      }

      setDropdownStyle({
        position: 'absolute',
        top: `${top}px`,
        left: `${left}px`,
        zIndex: 1000,
      });
    }
  }, [isOpen]);

  const dropdownContent = isOpen ? (
    <div
      ref={dropdownRef}
      className="w-48 bg-white rounded-md shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none"
      style={dropdownStyle}
    >
      <div className="py-1">
        {actions.map((action, index) => (
          <button
            key={index}
            onClick={(e) => {
              action.onClick(e);
              setIsOpen(false);
            }}
            className="flex items-center w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900"
            role="menuitem"
          >
            {action.icon && React.createElement(action.icon, { className: 'mr-3 h-4 w-4' })}
            {action.label}
          </button>
        ))}
      </div>
    </div>
  ) : null;

  return (
    <div className="relative">
      <Button
        variant="outline"
        size="sm"
        onClick={toggleDropdown}
        className="p-2 hover:bg-gray-100 hover:border-gray-300"
        title="More actions"
        ref={buttonRef}
      >
        <MoreVertical className="h-4 w-4 text-gray-600" />
      </Button>
      {isOpen && createPortal(dropdownContent, document.body)}
    </div>
  );
};

export default ActionDropdown;
