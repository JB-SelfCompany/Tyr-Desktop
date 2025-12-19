import React from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { pageTransitionVariants } from '../../styles/animations';

interface PageTransitionProps {
  children: React.ReactNode;
  pageKey?: string;
  className?: string;
}

export const PageTransition: React.FC<PageTransitionProps> = ({
  children,
  pageKey,
  className = '',
}) => {
  return (
    <AnimatePresence mode="wait">
      <motion.div
        key={pageKey}
        variants={pageTransitionVariants}
        initial="initial"
        animate="animate"
        exit="exit"
        className={className}
      >
        {children}
      </motion.div>
    </AnimatePresence>
  );
};
