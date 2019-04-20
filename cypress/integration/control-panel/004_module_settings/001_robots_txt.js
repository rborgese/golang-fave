/// <reference types="Cypress" />

context('Module robots.txt', () => {
  it('should reset', () => {
    cy.installCMS();
  });

  it('should render edit form', () => {
    cy.loginCMS();
    cy.visitCMS('/cp/settings/');
    cy.get('.data-form.settings- textarea[name=content]').should('exist');
    cy.get('.data-form.settings- textarea[name=content]').should('have.value', 'User-agent: *\nDisallow: /\n');
    cy.logoutCMS();
  });

  it('should render result file', () => {
    cy.request({
      url: cy.getBaseUrl() + '/robots.txt',
      followRedirect: false
    }).then((response) => {
      expect(response.status).to.eq(200);
      expect(response.body).to.eq('User-agent: *\r\nDisallow: /\r\n');
    });
  });

  it('should change file content', () => {
    cy.loginCMS();

    cy.visitCMS('/cp/settings/');
    cy.get('.data-form.settings- textarea[name=content]').clear().type('Some file content');
    cy.get('#add-edit-button').click();
    cy.actionWait();

    cy.visitCMS('/cp/settings/');
    cy.get('.data-form.settings- textarea[name=content]').should('have.value', 'Some file content');

    cy.request({
      url: cy.getBaseUrl() + '/robots.txt',
      followRedirect: false
    }).then((response) => {
      expect(response.status).to.eq(200);
      expect(response.body).to.eq('Some file content');
    });

    cy.visitCMS('/cp/settings/');
    cy.get('.data-form.settings- textarea[name=content]').clear().type('User-agent: *\nDisallow: /\n');
    cy.get('#add-edit-button').click();
    cy.actionWait();

    cy.visitCMS('/cp/settings/');
    cy.get('.data-form.settings- textarea[name=content]').should('have.value', 'User-agent: *\nDisallow: /\n');

    cy.request({
      url: cy.getBaseUrl() + '/robots.txt',
      followRedirect: false
    }).then((response) => {
      expect(response.status).to.eq(200);
      expect(response.body).to.eq('User-agent: *\r\nDisallow: /\r\n');
    });

    cy.logoutCMS();
  });
});