function errorHandler(err, req, res, next) {
  if (err && err.name === 'ValidationError') {
    return res.status(400).json({ error: err.message });
  }
  if (err && err.name === 'CastError') {
    return res.status(400).json({ error: 'Invalid id format' });
  }
  if (err && err.code === 11000) {
    return res.status(409).json({ error: 'Duplicate key error' });
  }
  console.error(err);
  return res.status(500).json({ error: 'Internal server error' });
}

module.exports = errorHandler;
