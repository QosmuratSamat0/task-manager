function toUserResponse(user) {
  if (!user) return null;
  return {
    id: String(user._id),
    _id: String(user._id),
    user_name: user.userName,
    userName: user.userName,
    email: user.email,
    role: user.role,
    created_at: user.createdAt,
    createdAt: user.createdAt,
    updated_at: user.updatedAt,
    updatedAt: user.updatedAt,
  };
}

function toProjectResponse(project) {
  if (!project) return null;
  return {
    id: String(project._id),
    _id: String(project._id),
    name: project.name,
    description: project.description,
    created_at: project.createdAt,
    createdAt: project.createdAt,
    updated_at: project.updatedAt,
    updatedAt: project.updatedAt,
  };
}

function toTaskResponse(task) {
  if (!task) return null;
  const userId = task.user && (task.user._id ? task.user._id : task.user);
  const projectId = task.project && (task.project._id ? task.project._id : task.project);
  const projectObj = task.project && task.project._id ? { id: String(task.project._id), name: task.project.name } : null;
  return {
    id: String(task._id),
    _id: String(task._id),
    user_id: userId ? String(userId) : null,
    project_id: projectId ? String(projectId) : null,
    project: projectObj,
    category: task.category || '',
    title: task.title,
    description: task.description,
    status: task.status,
    priority: task.priority,
    deadline: task.deadline,
    created_at: task.createdAt,
    updated_at: task.updatedAt,
  };
}

function toCategoryResponse(category) {
  if (!category) return null;
  return {
    id: String(category._id),
    name: category.name,
    description: category.description,
    color: category.color,
    created_at: category.createdAt,
    updated_at: category.updatedAt,
  };
}

module.exports = { toUserResponse, toProjectResponse, toTaskResponse, toCategoryResponse };
