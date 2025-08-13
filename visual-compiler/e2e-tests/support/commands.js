// ***********************************************
// This example commands.js shows you how to
// create various custom commands and overwrite
// existing commands.
//
// For more comprehensive examples of custom
// commands please read more here:
// https://on.cypress.io/custom-commands
// ***********************************************
//
//
// -- This is a parent command --
// Cypress.Commands.add('login', (email, password) => { ... })
//
//
// -- This is a child command --
// Cypress.Commands.add('drag', { prevSubject: 'element'}, (subject, options) => { ... })
//
//
// -- This is a dual command --
// Cypress.Commands.add('dismiss', { prevSubject: 'optional'}, (subject, options) => { ... })
//
//
// -- This will overwrite an existing command --
// Cypress.Commands.overwrite('visit', (originalFn, url, options) => { ... })

Cypress.Commands.add('ConnectNodes', (sourceNodeId, targetNodeId, sourceHandle = 'right', targetHandle = 'left') => {
  // Get the source node and its handle
    cy.get(`g[id*="${sourceNodeId}"]`).as('sourceNode');
    cy.get(`[id="${targetNodeId}"]`).as('targetNode');
  
  // Get the handle positions
  cy.get('@sourceNode').find(`.handle-${sourceHandle}`).as('sourceHandle');
  cy.get('@targetNode').find(`.handle-${targetHandle}`).as('targetHandle');
  
  // Perform the drag and drop to connect nodes
  cy.get('@sourceHandle').then(($sourceHandle) => {
    cy.get('@targetHandle').then(($targetHandle) => {
      const sourceRect = $sourceHandle[0].getBoundingClientRect();
      const targetRect = $targetHandle[0].getBoundingClientRect();
      
      // Start the drag from the source handle
      cy.get('@sourceHandle')
        .trigger('mousedown', {
          button: 0,
          clientX: sourceRect.left + sourceRect.width / 2,
          clientY: sourceRect.top + sourceRect.height / 2
        });
      
      // Move to the target handle
      cy.get('@targetHandle')
        .trigger('mousemove', {
          clientX: targetRect.left + targetRect.width / 2,
          clientY: targetRect.top + targetRect.height / 2
        });
      
      // Release to create the connection
      cy.get('@targetHandle')
        .trigger('mouseup', {
          clientX: targetRect.left + targetRect.width / 2,
          clientY: targetRect.top + targetRect.height / 2
        });
    });
  });
});

// Alternative approach using drag and drop with coordinates
Cypress.Commands.add('connectNodesByCoordinates', (sourceNodeId, targetNodeId) => {
  cy.get(`[data-node-id="${sourceNodeId}"]`).then($sourceNode => {
    cy.get(`[data-node-id="${targetNodeId}"]`).then($targetNode => {
      const sourceRect = $sourceNode[0].getBoundingClientRect();
      const targetRect = $targetNode[0].getBoundingClientRect();
      
      // Calculate handle positions (adjust based on your node structure)
      const sourceX = sourceRect.right - 10; // Right handle of source
      const sourceY = sourceRect.top + sourceRect.height / 2;
      const targetX = targetRect.left + 10; // Left handle of target
      const targetY = targetRect.top + targetRect.height / 2;
      
      // Perform the drag and drop
      cy.get(`[data-node-id="${sourceNodeId}"]`)
        .trigger('mousedown', { which: 1, clientX: sourceX, clientY: sourceY })
        .trigger('mousemove', { clientX: targetX, clientY: targetY })
        .trigger('mouseup', { clientX: targetX, clientY: targetY });
    });
  });
});

// Command to verify a connection exists between two nodes
Cypress.Commands.add('verifyNodeConnection', (sourceNodeId, targetNodeId) => {
  // Check if an edge exists between the nodes
  cy.get('.svelvet-edge').should('exist');
  
  // More specific verification - check for edge with specific source/target
  cy.get(`[data-edge-source="${sourceNodeId}"][data-edge-target="${targetNodeId}"]`)
    .should('exist');
});

// Command to disconnect nodes (if your implementation supports it)
Cypress.Commands.add('disconnectNodes', (sourceNodeId, targetNodeId) => {
  // Find the edge and click on it (if edges are clickable for deletion)
  cy.get(`[data-edge-source="${sourceNodeId}"][data-edge-target="${targetNodeId}"]`)
    .click()
    .type('{del}'); // or whatever key/action deletes edges in your app
});