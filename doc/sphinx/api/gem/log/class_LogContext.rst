LogContext
----------

.. py:currentmodule:: gem.log

.. py:class:: LogContext

   LogContext is a base class for objects which provide context to a logger.

   .. note:: Should not be used directly. Instead, subclass and provide implementations for methods.

   .. rubric:: Methods

   .. py:method:: log_context()

      Accessor for the context dictionary by this LogContext

      :return: The log context
      :rtype: dictionary
