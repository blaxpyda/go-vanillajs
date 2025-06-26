// Global state
let allProperties = [];
let allAgents = [];
let houseTypes = [];
let filteredProperties = [];

// DOM Content Loaded
document.addEventListener('DOMContentLoaded', function() {
    initializeApp();
    setupEventListeners();
});

// Initialize the application
function initializeApp() {
    loadTopProperties();
    loadAllProperties();
    loadAgents();
    loadHouseTypes();
}

// Setup event listeners
function setupEventListeners() {
    // Navigation
    document.querySelectorAll('.nav-link').forEach(link => {
        link.addEventListener('click', function(e) {
            e.preventDefault();
            const section = this.getAttribute('data-section');
            showSection(section);
            
            // Update active nav link
            document.querySelectorAll('.nav-link').forEach(l => l.classList.remove('active'));
            this.classList.add('active');
        });
    });

    // Filters
    const typeFilter = document.getElementById('type-filter');
    const priceFilter = document.getElementById('price-filter');
    
    if (typeFilter) {
        typeFilter.addEventListener('change', applyFilters);
    }
    
    if (priceFilter) {
        priceFilter.addEventListener('change', applyFilters);
    }

    // Modal close
    window.addEventListener('click', function(e) {
        const modal = document.getElementById('property-modal');
        if (e.target === modal) {
            closeModal();
        }
    });
}

// Show specific section
function showSection(sectionId) {
    // Hide all sections
    document.querySelectorAll('.section').forEach(section => {
        section.classList.remove('active');
    });
    
    // Show target section
    const targetSection = document.getElementById(sectionId);
    if (targetSection) {
        targetSection.classList.add('active');
    }
}

// API Functions
async function fetchAPI(endpoint) {
    try {
        const response = await fetch(endpoint);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        return await response.json();
    } catch (error) {
        console.error(`Error fetching ${endpoint}:`, error);
        return null;
    }
}

// Load top properties
async function loadTopProperties() {
    const loading = document.getElementById('home-loading');
    const container = document.getElementById('top-properties');
    
    if (loading) loading.style.display = 'block';
    
    try {
        const properties = await fetchAPI('/api/houses/top?limit=6');
        if (properties && container) {
            renderProperties(properties, container);
        }
    } catch (error) {
        console.error('Error loading top properties:', error);
        if (container) {
            container.innerHTML = '<p class="error">Failed to load top properties</p>';
        }
    } finally {
        if (loading) loading.style.display = 'none';
    }
}

// Load all properties
async function loadAllProperties() {
    const loading = document.getElementById('properties-loading');
    const container = document.getElementById('all-properties');
    
    if (loading) loading.style.display = 'block';
    
    try {
        const properties = await fetchAPI('/api/houses');
        if (properties) {
            allProperties = properties;
            filteredProperties = [...properties];
            if (container) {
                renderProperties(filteredProperties, container);
            }
        }
    } catch (error) {
        console.error('Error loading properties:', error);
        if (container) {
            container.innerHTML = '<p class="error">Failed to load properties</p>';
        }
    } finally {
        if (loading) loading.style.display = 'none';
    }
}

// Load agents
async function loadAgents() {
    const loading = document.getElementById('agents-loading');
    const container = document.getElementById('agents-list');
    
    if (loading) loading.style.display = 'block';
    
    try {
        const agents = await fetchAPI('/api/agents');
        if (agents) {
            allAgents = agents;
            if (container) {
                renderAgents(agents, container);
            }
        }
    } catch (error) {
        console.error('Error loading agents:', error);
        if (container) {
            container.innerHTML = '<p class="error">Failed to load agents</p>';
        }
    } finally {
        if (loading) loading.style.display = 'none';
    }
}

// Load house types
async function loadHouseTypes() {
    try {
        const types = await fetchAPI('/api/house-types');
        if (types) {
            houseTypes = types;
            populateTypeFilter(types);
        }
    } catch (error) {
        console.error('Error loading house types:', error);
    }
}

// Render properties
function renderProperties(properties, container) {
    if (!container) return;
    
    if (!properties || properties.length === 0) {
        container.innerHTML = '<p class="no-results">No properties found</p>';
        return;
    }
    
    container.innerHTML = properties.map(property => `
        <div class="property-card" onclick="showPropertyDetails(${property.id})">
            <img src="${property.image_url || '/images/logo.png'}" alt="${property.name}" class="property-image" onerror="this.src='/images/logo.png'">
            <div class="property-info">
                <h3 class="property-name">${property.name}</h3>
                <p class="property-description">${truncateText(property.description, 100)}</p>
                <div class="property-price">$${formatPrice(property.price)}</div>
                <div class="property-details">
                    <div class="property-type">${property.house_type ? property.house_type.name : 'Unknown'}</div>
                    ${property.agent ? `
                        <div class="property-agent">
                            <img src="${property.agent.image_url || '/images/generic_actor.jpg'}" alt="${property.agent.first_name}" class="agent-avatar" onerror="this.src='/images/generic_actor.jpg'">
                            <span>${property.agent.first_name} ${property.agent.last_name}</span>
                        </div>
                    ` : ''}
                </div>
                ${property.tags && property.tags.length > 0 ? `
                    <div class="property-tags">
                        ${property.tags.map(tag => `<span class="tag">${tag}</span>`).join('')}
                    </div>
                ` : ''}
            </div>
        </div>
    `).join('');
}

// Render agents
function renderAgents(agents, container) {
    if (!container) return;
    
    if (!agents || agents.length === 0) {
        container.innerHTML = '<p class="no-results">No agents found</p>';
        return;
    }
    
    container.innerHTML = agents.map(agent => `
        <div class="agent-card">
            <img src="${agent.image_url || '/images/generic_actor.jpg'}" alt="${agent.first_name}" class="agent-image" onerror="this.src='/images/generic_actor.jpg'">
            <h3 class="agent-name">${agent.first_name} ${agent.last_name}</h3>
        </div>
    `).join('');
}

// Populate type filter
function populateTypeFilter(types) {
    const filter = document.getElementById('type-filter');
    if (!filter) return;
    
    types.forEach(type => {
        const option = document.createElement('option');
        option.value = type.id;
        option.textContent = type.name;
        filter.appendChild(option);
    });
}

// Apply filters
function applyFilters() {
    const typeFilter = document.getElementById('type-filter');
    const priceFilter = document.getElementById('price-filter');
    
    let filtered = [...allProperties];
    
    // Filter by type
    if (typeFilter && typeFilter.value) {
        const typeId = parseInt(typeFilter.value);
        filtered = filtered.filter(property => property.house_type_id === typeId);
    }
    
    // Filter by price
    if (priceFilter && priceFilter.value) {
        const priceRange = priceFilter.value;
        filtered = filtered.filter(property => {
            const price = property.price;
            switch (priceRange) {
                case '0-300000':
                    return price < 300000;
                case '300000-600000':
                    return price >= 300000 && price < 600000;
                case '600000-1000000':
                    return price >= 600000 && price < 1000000;
                case '1000000+':
                    return price >= 1000000;
                default:
                    return true;
            }
        });
    }
    
    filteredProperties = filtered;
    const container = document.getElementById('all-properties');
    renderProperties(filteredProperties, container);
}

// Show property details in modal
async function showPropertyDetails(propertyId) {
    const modal = document.getElementById('property-modal');
    const modalBody = document.getElementById('modal-body');
    
    if (!modal || !modalBody) return;
    
    modalBody.innerHTML = '<div class="loading">Loading property details...</div>';
    modal.style.display = 'block';
    
    try {
        const property = await fetchAPI(`/api/houses/${propertyId}`);
        if (property) {
            modalBody.innerHTML = `
                <div class="property-details-modal">
                    <img src="${property.image_url || '/images/logo.png'}" alt="${property.name}" class="property-image-large" onerror="this.src='/images/logo.png'">
                    <h2>${property.name}</h2>
                    <div class="property-price-large">$${formatPrice(property.price)}</div>
                    <p class="property-description-full">${property.description}</p>
                    
                    <div class="property-info-grid">
                        <div class="info-item">
                            <strong>Type:</strong> ${property.house_type ? property.house_type.name : 'Unknown'}
                        </div>
                        ${property.agent ? `
                            <div class="info-item">
                                <strong>Agent:</strong> 
                                <div class="agent-info">
                                    <img src="${property.agent.image_url || '/images/generic_actor.jpg'}" alt="${property.agent.first_name}" class="agent-avatar" onerror="this.src='/images/generic_actor.jpg'">
                                    <span>${property.agent.first_name} ${property.agent.last_name}</span>
                                </div>
                            </div>
                        ` : ''}
                        <div class="info-item">
                            <strong>Listed:</strong> ${formatDate(property.created_at)}
                        </div>
                        <div class="info-item">
                            <strong>Updated:</strong> ${formatDate(property.updated_at)}
                        </div>
                    </div>
                    
                    ${property.tags && property.tags.length > 0 ? `
                        <div class="property-tags">
                            <strong>Tags:</strong>
                            ${property.tags.map(tag => `<span class="tag">${tag}</span>`).join('')}
                        </div>
                    ` : ''}
                </div>
            `;
        }
    } catch (error) {
        console.error('Error loading property details:', error);
        modalBody.innerHTML = '<p class="error">Failed to load property details</p>';
    }
}

// Close modal
function closeModal() {
    const modal = document.getElementById('property-modal');
    if (modal) {
        modal.style.display = 'none';
    }
}

// Utility functions
function formatPrice(price) {
    return new Intl.NumberFormat('en-US').format(price);
}

function truncateText(text, maxLength) {
    if (!text) return '';
    if (text.length <= maxLength) return text;
    return text.substring(0, maxLength) + '...';
}

function formatDate(dateString) {
    if (!dateString) return 'Unknown';
    try {
        const date = new Date(dateString);
        return date.toLocaleDateString('en-US', { 
            year: 'numeric', 
            month: 'short', 
            day: 'numeric' 
        });
    } catch (error) {
        return 'Unknown';
    }
}

// Add some additional CSS for modal content
const modalStyles = `
    .property-details-modal {
        text-align: center;
    }
    
    .property-image-large {
        width: 100%;
        height: 300px;
        object-fit: cover;
        border-radius: 8px;
        margin-bottom: 1rem;
    }
    
    .property-price-large {
        font-size: 2rem;
        font-weight: bold;
        color: #667eea;
        margin: 1rem 0;
    }
    
    .property-description-full {
        text-align: left;
        margin: 1.5rem 0;
        line-height: 1.6;
    }
    
    .property-info-grid {
        display: grid;
        gap: 1rem;
        text-align: left;
        margin: 1.5rem 0;
    }
    
    .info-item {
        display: flex;
        align-items: center;
        gap: 0.5rem;
    }
    
    .agent-info {
        display: flex;
        align-items: center;
        gap: 0.5rem;
    }
    
    .error {
        color: #e53e3e;
        text-align: center;
        padding: 2rem;
    }
    
    .no-results {
        text-align: center;
        padding: 2rem;
        color: #718096;
        font-style: italic;
    }
`;

// Inject additional styles
const styleSheet = document.createElement('style');
styleSheet.textContent = modalStyles;
document.head.appendChild(styleSheet);
