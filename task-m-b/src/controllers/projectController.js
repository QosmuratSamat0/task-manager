const Project = require('../models/Project');
const { toProjectResponse } = require('../utils/formatters');

async function getProjects(req, res) {
  const projects = await Project.find().sort({ createdAt: -1 });
  return res.json({ data: projects.map(toProjectResponse) });
}

async function getProjectById(req, res) {
  const project = await Project.findById(req.params.id);
  if (!project) {
    return res.status(404).json({ error: 'Project not found' });
  }
  return res.json({ data: toProjectResponse(project) });
}

async function createProject(req, res) {
  const project = await Project.create(req.body);
  return res.status(201).json({ data: toProjectResponse(project) });
}

async function updateProject(req, res) {
  const project = await Project.findByIdAndUpdate(req.params.id, req.body, {
    new: true,
    runValidators: true,
  });
  if (!project) {
    return res.status(404).json({ error: 'Project not found' });
  }
  return res.json({ data: toProjectResponse(project) });
}

async function deleteProject(req, res) {
  const project = await Project.findByIdAndDelete(req.params.id);
  if (!project) {
    return res.status(404).json({ error: 'Project not found' });
  }
  return res.json({ success: true });
}

module.exports = {
  getProjects,
  getProjectById,
  createProject,
  updateProject,
  deleteProject,
};
